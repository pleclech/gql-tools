package gqlexpr

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vektah/gqlparser/v2/ast"
)

// OperatorCatalog is the minimal operator metadata needed for parsing.
type OperatorCatalog interface {
	OperatorByDirective(directive string) (name string, sql string, args []ArgSpec, ok bool)
}

// ArgSpec describes one operator argument.
type ArgSpec struct {
	Name     string
	Kind     string // expr | expr_list
	MinItems int
}

// CatalogAdapter wraps a simple map for tests and exprconfig.
type CatalogAdapter struct {
	ByDirective map[string]struct {
		Name string
		SQL  string
		Args []ArgSpec
	}
}

func (c *CatalogAdapter) OperatorByDirective(directive string) (string, string, []ArgSpec, bool) {
	if c == nil {
		return "", "", nil, false
	}
	op, ok := c.ByDirective[directive]
	if !ok {
		return "", "", nil, false
	}
	return op.Name, op.SQL, op.Args, true
}

// ParseDirective converts an @x_* expression directive to a Node AST.
func ParseDirective(d *ast.Directive, cat OperatorCatalog, tableContext string, cols *ColumnRegistry) (Node, error) {
	return parseOperatorArg(d.Name, d.Name, d.Arguments, cat, tableContext, cols)
}

// ParseFieldArgument converts a field argument (e.g. x_concat: { parts: [...] }) to a Node AST.
func ParseFieldArgument(arg *ast.Argument, cat OperatorCatalog, tableContext string, cols *ColumnRegistry) (Node, error) {
	if arg == nil {
		return Node{}, fmt.Errorf("expression argument is nil")
	}
	if arg.Value == nil {
		return Node{}, fmt.Errorf("expression argument %q: missing value", arg.Name)
	}
	return parseOperatorArgValue(arg.Name, arg.Name, arg.Value, cat, tableContext, cols)
}

func parseOperatorArg(directiveName, label string, args ast.ArgumentList, cat OperatorCatalog, tableContext string, cols *ColumnRegistry) (Node, error) {
	_, _, argSpecs, ok := cat.OperatorByDirective(directiveName)
	if !ok {
		return Node{}, fmt.Errorf("unknown expression operator %q", label)
	}
	if len(argSpecs) != 1 {
		return Node{}, fmt.Errorf("operator %q: expected one argument definition", label)
	}
	spec := argSpecs[0]
	var argVal *ast.Value
	for _, a := range args {
		if a.Name == spec.Name {
			argVal = a.Value
			break
		}
	}
	if argVal == nil {
		return Node{}, fmt.Errorf("operator %q: missing argument %q", label, spec.Name)
	}
	return parseOperatorArgValue(directiveName, label, argVal, cat, tableContext, cols)
}

func parseOperatorArgValue(directiveName, label string, argVal *ast.Value, cat OperatorCatalog, tableContext string, cols *ColumnRegistry) (Node, error) {
	opName, _, argSpecs, ok := cat.OperatorByDirective(directiveName)
	if !ok {
		return Node{}, fmt.Errorf("unknown expression operator %q", label)
	}
	spec := argSpecs[0]
	// Field args use an input object wrapper: x_concat: { parts: [...] }.
	if argVal.Kind == ast.ObjectValue {
		if inner := objectField(argVal, spec.Name); inner != nil {
			argVal = inner
		}
	}
	var nodes []Node
	var err error
	if spec.Kind == "expr_list" {
		nodes, err = parseExprListValue(argVal, tableContext, cols)
	} else {
		var n Node
		n, err = parseExprValue(argVal, tableContext, cols)
		nodes = []Node{n}
	}
	if err != nil {
		return Node{}, err
	}
	if spec.MinItems > 0 && len(nodes) < spec.MinItems {
		return Node{}, fmt.Errorf("operator %q: need at least %d expressions", label, spec.MinItems)
	}
	return Node{Op: opName, Table: tableContext, Args: nodes}, nil
}

func parseExprListValue(v *ast.Value, tableContext string, cols *ColumnRegistry) ([]Node, error) {
	if v == nil || v.Kind != ast.ListValue {
		return nil, fmt.Errorf("expected list value")
	}
	out := make([]Node, 0, len(v.Children))
	for _, child := range v.Children {
		if child.Value == nil {
			continue
		}
		n, err := parseExprValue(child.Value, tableContext, cols)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}

func parseExprValue(v *ast.Value, tableContext string, cols *ColumnRegistry) (Node, error) {
	if v == nil || v.Kind != ast.ObjectValue {
		return Node{}, fmt.Errorf("expected object value")
	}
	// column: { name: ENUM } or shorthand
	if colField := objectField(v, "column"); colField != nil {
		name, err := columnNameFromValue(colField)
		if err != nil {
			return Node{}, err
		}
		if cols != nil && !cols.Allowed(tableContext, name) {
			return Node{}, fmt.Errorf("column %q is not allowed on %q", name, tableContext)
		}
		return Node{Col: name}, nil
	}
	if strField := objectField(v, "string"); strField != nil {
		s, err := stringLiteralFromValue(strField)
		if err != nil {
			return Node{}, err
		}
		return Node{Str: s}, nil
	}
	// nested operator: concat: { parts: [...] }
	for _, c := range v.Children {
		if c.Value == nil || c.Value.Kind != ast.ObjectValue {
			continue
		}
		// inline nested op not used in v1 client syntax (directives only)
		_ = c
	}
	// Allow nested op keys matching catalog would need catalog here — skip for directive-only v1
	return Node{}, fmt.Errorf("expression object must contain column or string")
}

func objectField(v *ast.Value, name string) *ast.Value {
	for _, c := range v.Children {
		if c.Name == name {
			return c.Value
		}
	}
	return nil
}

func columnNameFromValue(v *ast.Value) (string, error) {
	if v == nil {
		return "", fmt.Errorf("column value is nil")
	}
	// Shorthand: { column: firstname } (enum) or legacy { column: { name: firstname } }.
	if v.Kind == ast.ObjectValue {
		nf := objectField(v, "name")
		if nf == nil {
			return "", fmt.Errorf("column.name required")
		}
		v = nf
	}
	if v.Kind != ast.StringValue && v.Kind != ast.EnumValue {
		return "", fmt.Errorf("column name must be string/enum")
	}
	s, err := v.Value(nil)
	if err != nil {
		return "", err
	}
	return s.(string), nil
}

func stringLiteralFromValue(v *ast.Value) (string, error) {
	if v == nil {
		return "", fmt.Errorf("string value is nil")
	}
	// Shorthand: { string: " " } or legacy { string: { value: " " } }.
	if v.Kind == ast.ObjectValue {
		valField := objectField(v, "value")
		if valField == nil || valField.Kind != ast.StringValue {
			return "", fmt.Errorf("string.value required")
		}
		return decodeGraphQLString(valField)
	}
	return decodeGraphQLString(v)
}

func decodeGraphQLString(v *ast.Value) (string, error) {
	if v == nil || v.Kind != ast.StringValue {
		return "", fmt.Errorf("string literal must be a string")
	}
	if u, err := strconv.Unquote(v.Raw); err == nil {
		return u, nil
	}
	s, err := v.Value(nil)
	if err != nil {
		return "", err
	}
	out, ok := s.(string)
	if !ok {
		return "", fmt.Errorf("string literal: invalid value")
	}
	return out, nil
}

// DirectiveFromCatalog builds a CatalogAdapter from exprconfig-like maps.
func DirectiveFromCatalog(byDirective map[string]struct {
	Name string
	SQL  string
	Args []ArgSpec
}) OperatorCatalog {
	return &CatalogAdapter{ByDirective: byDirective}
}

// ExprToRawMessage marshals a Node for NamedObject.ex.
func ExprToRawMessage(n Node) (json.RawMessage, error) {
	return json.Marshal(n)
}
