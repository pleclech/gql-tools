package gqlhasura

import (
	"testing"

	"github.com/pleclech/gql-tools/gqldoc"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestNormalizeOrToList(t *testing.T) {
	where, err := gqldoc.ParseWhereLiteral(`_or:{email:{_eq:"a"},phone:{_eq:"b"}}`)
	if err != nil {
		t.Fatalf("ParseWhereLiteral: %v", err)
	}

	ok, nv := NormalizeOrToList(where)
	if !ok {
		t.Fatal("expected modification")
	}
	if nv.Value.Kind != ast.ListValue {
		t.Fatalf("expected list value, got kind %v", nv.Value.Kind)
	}
	if len(nv.Value.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(nv.Value.Children))
	}

	ok, _ = NormalizeOrToList(nv)
	if ok {
		t.Fatal("expected no-op on already-list _or")
	}
}

func TestNormalizeWhereOrLists(t *testing.T) {
	where, err := gqldoc.ParseValueLiteral(`{_and:[{_or:{a:{_eq:1},b:{_eq:2}}}]} `)
	if err != nil {
		t.Fatalf("ParseValueLiteral: %v", err)
	}

	if !NormalizeWhereOrLists(where) {
		t.Fatal("expected nested _or to be normalized")
	}
}
