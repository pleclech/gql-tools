package gqldoc

import (
	"github.com/vektah/gqlparser/v2/ast"
)

// WalkOperations calls fn for each operation in the document.
func WalkOperations(doc *ast.QueryDocument, fn func(*ast.OperationDefinition) error) error {
	if doc == nil {
		return nil
	}
	for _, op := range doc.Operations {
		if op == nil {
			continue
		}
		if err := fn(op); err != nil {
			return err
		}
	}
	return nil
}

// WalkFields calls fn for each field in the selection set.
func WalkFields(ss ast.SelectionSet, fn func(*ast.Field) error) error {
	for _, sel := range ss {
		field, ok := AsField(sel)
		if !ok {
			continue
		}
		if err := fn(field); err != nil {
			return err
		}
	}
	return nil
}

// WalkFieldArguments calls fn for each argument on a field.
func WalkFieldArguments(field *ast.Field, fn func(*ast.Argument) error) error {
	if field == nil {
		return nil
	}
	for _, arg := range field.Arguments {
		if arg == nil {
			continue
		}
		if err := fn(arg); err != nil {
			return err
		}
	}
	return nil
}

// MapWhereArguments calls fn for each root field that has a where object argument.
func MapWhereArguments(op *ast.OperationDefinition, fn func(field *ast.Field, where *ast.Value) error) error {
	if op == nil {
		return nil
	}
	return WalkFields(op.SelectionSet, func(field *ast.Field) error {
		arg, ok := FindWhere(field)
		if !ok {
			return nil
		}
		return fn(field, arg.Value)
	})
}
