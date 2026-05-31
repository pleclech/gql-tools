package gqlhasura

import (
	"github.com/pleclech/gql-tools/gqldoc"
	"github.com/vektah/gqlparser/v2/ast"
)

// NormalizeOrToList converts Hasura-style {_or:{a,b}} object filters to {_or:[{a},{b}]}.
// Returns (true, node) when modified, (false, node) when unchanged.
func NormalizeOrToList(op *ast.ChildValue) (bool, *ast.ChildValue) {
	if op == nil || op.Name != "_or" {
		return false, op
	}

	value := op.Value
	if value == nil || value.Kind == ast.ListValue {
		return false, op
	}

	children := make(ast.ChildValueList, 0, len(value.Children))
	for _, v := range value.Children {
		children = append(children, &ast.ChildValue{
			Value: &ast.Value{
				Kind:     ast.ObjectValue,
				Children: []*ast.ChildValue{v},
			},
		})
	}

	op.Value.Kind = ast.ListValue
	op.Value.Children = children
	return true, op
}

// NormalizeWhereOrLists walks a where object value and normalizes every _or field to list form.
// Returns true when at least one node was modified.
func NormalizeWhereOrLists(where *ast.Value) bool {
	if where == nil || where.Kind != ast.ObjectValue {
		return false
	}

	modified := false
	gqldoc.MapObjectFields(where, func(path gqldoc.ObjectPath) (*ast.ChildValue, bool) {
		if ok, nv := NormalizeOrToList(path.Node); ok {
			modified = true
			return nv, true
		}
		return nil, false
	})
	return modified
}
