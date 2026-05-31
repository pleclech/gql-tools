package gqldoc

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestAccessHelpers(t *testing.T) {
	doc, err := ParseQuery(`query Q($id: Int!) { users(where: { id: { _eq: $id } }, limit: 10) { id name } }`)
	if err != nil {
		t.Fatal(err)
	}
	op, err := FirstOperation(doc)
	if err != nil {
		t.Fatal(err)
	}
	field, ok := FirstField(op.SelectionSet)
	if !ok {
		t.Fatal("expected field")
	}
	if FieldResponseName(field) != "users" {
		t.Fatalf("got %q", FieldResponseName(field))
	}
	arg, ok := FindWhere(field)
	if !ok {
		t.Fatal("expected where")
	}
	if arg.Value.Kind != ast.ObjectValue {
		t.Fatalf("expected object value")
	}
}

func TestVariableName(t *testing.T) {
	if VariableName("$foo") != "foo" {
		t.Fatal("expected foo")
	}
}
