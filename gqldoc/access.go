package gqldoc

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// AsField returns the selection as a field when possible.
func AsField(sel ast.Selection) (*ast.Field, bool) {
	field, ok := sel.(*ast.Field)
	return field, ok
}

// FirstOperation returns the first operation in a document.
func FirstOperation(doc *ast.QueryDocument) (*ast.OperationDefinition, error) {
	if doc == nil || len(doc.Operations) == 0 {
		return nil, fmt.Errorf("no operations in document")
	}
	return doc.Operations[0], nil
}

// FirstField returns the first field in a selection set.
func FirstField(ss ast.SelectionSet) (*ast.Field, bool) {
	if len(ss) == 0 {
		return nil, false
	}
	return AsField(ss[0])
}

// FieldResponseName returns the alias if set, otherwise the field name.
func FieldResponseName(f *ast.Field) string {
	if f == nil {
		return ""
	}
	if f.Alias != "" {
		return f.Alias
	}
	return f.Name
}

// FindArgument returns an argument by name and its index.
func FindArgument(args ast.ArgumentList, name string) (*ast.Argument, int, bool) {
	for i, arg := range args {
		if arg != nil && arg.Name == name {
			return arg, i, true
		}
	}
	return nil, -1, false
}

// FindFieldArgument returns a field argument by name.
func FindFieldArgument(field *ast.Field, name string) (*ast.Argument, bool) {
	if field == nil {
		return nil, false
	}
	arg, _, ok := FindArgument(field.Arguments, name)
	return arg, ok
}

// FindWhere returns the where argument when it holds an object value.
func FindWhere(field *ast.Field) (*ast.Argument, bool) {
	arg, ok := FindFieldArgument(field, "where")
	if !ok || arg.Value == nil || arg.Value.Kind != ast.ObjectValue {
		return nil, false
	}
	return arg, true
}

// VariableName strips the leading $ from a variable reference raw value.
func VariableName(raw string) string {
	return strings.TrimPrefix(raw, "$")
}
