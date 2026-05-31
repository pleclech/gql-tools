package gqldoc

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestBuildHelpers(t *testing.T) {
	v := Obj(
		Field("status", Obj(
			Field("code", Obj(
				Field("_eq", Str("active")),
			)),
		)),
	)
	if v.Kind != ast.ObjectValue || len(v.Children) != 1 {
		t.Fatalf("unexpected object: %+v", v)
	}
	if Var("id").Raw != "$id" {
		t.Fatal("expected $id")
	}
	if Var("$id").Raw != "$id" {
		t.Fatal("expected $id")
	}
}

func TestWrapNot(t *testing.T) {
	inner := Field("_and", Obj(Field("a", Bool("true"))))
	wrapped := WrapNot(inner)
	if wrapped.Name != "_not" {
		t.Fatalf("got %q", wrapped.Name)
	}
	if len(wrapped.Value.Children) != 1 || wrapped.Value.Children[0] != inner {
		t.Fatalf("unexpected children: %+v", wrapped.Value.Children)
	}
}
