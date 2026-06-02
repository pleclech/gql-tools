package gqlexpr

import (
	"strings"
	"testing"
)

func TestBuildExtensionSDL(t *testing.T) {
	cat := &CatalogAdapter{
		ByDirective: map[string]struct {
			Name string
			SQL  string
			Args []ArgSpec
		}{
			"x_concat": {Name: "concat", SQL: "concat_ws", Args: []ArgSpec{{Name: "parts", Kind: "expr_list", MinItems: 2}}},
		},
	}
	enums := map[string][]string{
		"users_select_column": {"id", "firstname", "lastname"},
	}
	sdl := BuildExtensionSDL(cat, enums)
	if !strings.Contains(sdl, "column: users_select_column") {
		t.Fatalf("expected Hasura column enum on UsersExpr:\n%s", sdl)
	}
	if !strings.Contains(sdl, "string: String") {
		t.Fatalf("expected string shorthand on UsersExpr:\n%s", sdl)
	}
	if !strings.Contains(sdl, "directive @x_concat") {
		t.Fatalf("missing x_concat directive:\n%s", sdl)
	}
	if !strings.Contains(sdl, "input UsersExpr") {
		t.Fatalf("missing UsersExpr input:\n%s", sdl)
	}
}
