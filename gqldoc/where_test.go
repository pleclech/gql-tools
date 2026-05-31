package gqldoc

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestWhereHelpers(t *testing.T) {
	where := Obj(
		Field("status", Obj(
			Field("code", Obj(
				Field("_eq", Var("opt")),
			)),
		)),
		Field("_and", Obj(
			Field("a", Var("other")),
		)),
	)

	vars := CollectVariables(where)
	if len(vars["opt"]) != 1 {
		t.Fatalf("expected 1 opt path, got %d", len(vars["opt"]))
	}
	if len(vars["other"]) != 1 {
		t.Fatalf("expected 1 other path, got %d", len(vars["other"]))
	}

	replacement := Field("status", Obj())
	MapObjectFields(where, func(path ObjectPath) (*ast.ChildValue, bool) {
		if path.Node.Name == "status" {
			return replacement, true
		}
		return path.Node, false
	})
	if where.Children[0].Name != "status" || len(where.Children[0].Value.Children) != 0 {
		t.Fatalf("replacement failed: %+v", where.Children[0])
	}
}

func TestReplaceAt(t *testing.T) {
	where := Obj(
		Field("a", Obj(Field("b", Str("1")))),
	)
	paths := CollectVariables(where)
	_ = paths

	empty := Field("a", Obj())
	_ = WalkObjectFields(where, func(path ObjectPath) error {
		if path.Node.Name == "a" {
			ReplaceAt(path, empty)
		}
		return nil
	})
	if len(where.Children[0].Value.Children) != 0 {
		t.Fatalf("expected empty nested object")
	}
}
