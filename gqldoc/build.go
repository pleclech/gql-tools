package gqldoc

import (
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// Obj builds an object value from named fields.
func Obj(fields ...*ast.ChildValue) *ast.Value {
	return &ast.Value{
		Kind:     ast.ObjectValue,
		Children: fields,
	}
}

// Field builds a named object field.
func Field(name string, value *ast.Value) *ast.ChildValue {
	return &ast.ChildValue{
		Name:  name,
		Value: value,
	}
}

// Str builds a string value.
func Str(raw string) *ast.Value {
	return &ast.Value{Kind: ast.StringValue, Raw: raw}
}

// Int builds an integer value.
func Int(raw string) *ast.Value {
	return &ast.Value{Kind: ast.IntValue, Raw: raw}
}

// Bool builds a boolean value.
func Bool(raw string) *ast.Value {
	return &ast.Value{Kind: ast.BooleanValue, Raw: raw}
}

// Null builds a null value.
func Null() *ast.Value {
	return &ast.Value{Kind: ast.NullValue}
}

// Var builds a variable reference. The name may include or omit a leading $.
func Var(name string) *ast.Value {
	raw := name
	if !strings.HasPrefix(raw, "$") {
		raw = "$" + raw
	}
	return &ast.Value{Kind: ast.Variable, Raw: raw}
}

// List builds a list value from child values.
func List(values ...*ast.ChildValue) *ast.Value {
	return &ast.Value{
		Kind:     ast.ListValue,
		Children: values,
	}
}

// WrapNot wraps a child value in a _not object field.
func WrapNot(inner *ast.ChildValue) *ast.ChildValue {
	return &ast.ChildValue{
		Name: "_not",
		Value: &ast.Value{
			Kind:     ast.ObjectValue,
			Children: []*ast.ChildValue{inner},
		},
		Position: inner.Position,
		Comment:  inner.Comment,
	}
}

// CloneChild returns a shallow copy of a child value preserving metadata.
func CloneChild(cv *ast.ChildValue) *ast.ChildValue {
	if cv == nil {
		return nil
	}
	clone := *cv
	if cv.Value != nil {
		valueClone := *cv.Value
		if len(cv.Value.Children) > 0 {
			valueClone.Children = append(ast.ChildValueList(nil), cv.Value.Children...)
		}
		clone.Value = &valueClone
	}
	return &clone
}
