package gqldoc

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestFilterHelpers(t *testing.T) {
	args := ast.ArgumentList{
		{Name: "where", Value: Obj()},
		{Name: "limit", Value: Int("10")},
		{Name: "x_having", Value: Bool("true")},
	}
	filtered := RemoveArgumentsByName(args, "x_having")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 args, got %d", len(filtered))
	}

	dirs := ast.DirectiveList{
		{Name: "x_group_by"},
		{Name: "include"},
	}
	filteredDirs := RemoveDirectivesByPrefix(dirs, "x_")
	if len(filteredDirs) != 1 || filteredDirs[0].Name != "include" {
		t.Fatalf("unexpected directives: %+v", filteredDirs)
	}
}
