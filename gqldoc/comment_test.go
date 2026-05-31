package gqldoc

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestCommentHelpers(t *testing.T) {
	doc, err := ParseQuery(`query { f { id } }`)
	if err != nil {
		t.Fatal(err)
	}
	field, ok := FirstField(doc.Operations[0].SelectionSet)
	if !ok {
		t.Fatal("expected field")
	}
	field.Comment = &ast.CommentGroup{
		List: []*ast.Comment{{Value: "# @fallback: {a:{b:true}}"}},
	}
	payload, ok := ParseCommentPrefix(field.Comment, "@fallback:")
	if !ok {
		t.Fatal("expected fallback prefix")
	}
	if payload != "{a:{b:true}}" {
		t.Fatalf("got %q", payload)
	}
}
