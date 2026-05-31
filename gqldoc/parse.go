package gqldoc

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

// ParseQuery parses a GraphQL query document string.
func ParseQuery(input string) (*ast.QueryDocument, error) {
	return parser.ParseQuery(&ast.Source{Input: input})
}

// ParseWrapped parses fragment by wrapping it in a synthetic query operation.
// The fragment is inserted after the operation keyword/name and before the selection set.
// Example: ParseWrapped("q @d(w: true) {i}", "query") parses directives on the operation.
func ParseWrapped(name, fragment string) (*ast.QueryDocument, error) {
	query := "query " + name + " " + fragment + " {i}"
	return ParseQuery(query)
}

// ParseDirectives parses a directive list from a fragment placed before a dummy field.
func ParseDirectives(fragment string) ([]*ast.Directive, error) {
	doc, err := ParseWrapped("q", fragment)
	if err != nil {
		return nil, err
	}
	if len(doc.Operations) == 0 {
		return nil, nil
	}
	return doc.Operations[0].Directives, nil
}

// ParseDirective parses a single directive from a fragment.
func ParseDirective(fragment string) (*ast.Directive, error) {
	directives, err := ParseDirectives(fragment)
	if err != nil {
		return nil, err
	}
	if len(directives) == 0 {
		return nil, nil
	}
	return directives[0], nil
}

// ParseValueLiteral parses a GraphQL input value from a fragment.
func ParseValueLiteral(fragment string) (*ast.Value, error) {
	d, err := ParseDirective("@d(w:" + fragment + ")")
	if err != nil {
		return nil, err
	}
	if len(d.Arguments) == 0 {
		return nil, fmt.Errorf("invalid value literal: no argument")
	}
	return d.Arguments[0].Value, nil
}

// ParseWhereLiteral parses a Hasura-style where object field from a fragment.
func ParseWhereLiteral(fragment string) (*ast.ChildValue, error) {
	fragment = strings.TrimSpace(fragment)
	if fragment == "" {
		return nil, fmt.Errorf("where fragment is empty")
	}

	whereLiteral := fragment
	if !strings.HasPrefix(whereLiteral, "{") {
		whereLiteral = "{" + whereLiteral + "}"
	}

	queryDoc, err := ParseQuery("query __parseWhere { _x(where: " + whereLiteral + ") }")
	if err != nil {
		return nil, fmt.Errorf("failed to parse where fragment: %w", err)
	}
	if len(queryDoc.Operations) == 0 || len(queryDoc.Operations[0].SelectionSet) == 0 {
		return nil, fmt.Errorf("invalid where fragment: no operation or selection")
	}

	field, ok := AsField(queryDoc.Operations[0].SelectionSet[0])
	if !ok {
		return nil, fmt.Errorf("invalid where fragment: first selection is not a field")
	}

	arg, ok := FindFieldArgument(field, "where")
	if !ok || arg.Value == nil {
		return nil, fmt.Errorf("invalid where fragment: no where argument found")
	}
	if len(arg.Value.Children) != 1 {
		return nil, fmt.Errorf("invalid where fragment: where argument is not a single value")
	}
	return arg.Value.Children[0], nil
}
