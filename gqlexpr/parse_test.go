package gqlexpr

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestParseDirectiveConcat(t *testing.T) {
	cat := &CatalogAdapter{
		ByDirective: map[string]struct {
			Name string
			SQL  string
			Args []ArgSpec
		}{
			"x_concat": {Name: "concat", SQL: "concat_ws", Args: []ArgSpec{{Name: "parts", Kind: "expr_list", MinItems: 2}}},
		},
	}
	cols := NewColumnRegistry(map[string][]string{
		"users_select_column": {"firstname", "lastname"},
	})

	d := &ast.Directive{
		Name: "x_concat",
		Arguments: ast.ArgumentList{
			{
				Name: "parts",
				Value: &ast.Value{
					Kind: ast.ListValue,
					Children: ast.ChildValueList{
						{Value: &ast.Value{
							Kind: ast.ObjectValue,
							Children: ast.ChildValueList{
								{Name: "column", Value: &ast.Value{
									Kind: ast.ObjectValue,
									Children: ast.ChildValueList{
										{Name: "name", Value: &ast.Value{Kind: ast.StringValue, Raw: "firstname"}},
									},
								}},
							},
						}},
						{Value: &ast.Value{
							Kind: ast.ObjectValue,
							Children: ast.ChildValueList{
								{Name: "string", Value: &ast.Value{
									Kind: ast.ObjectValue,
									Children: ast.ChildValueList{
										{Name: "value", Value: &ast.Value{Kind: ast.StringValue, Raw: `" "`}},
									},
								}},
							},
						}},
						{Value: &ast.Value{
							Kind: ast.ObjectValue,
							Children: ast.ChildValueList{
								{Name: "column", Value: &ast.Value{
									Kind: ast.ObjectValue,
									Children: ast.ChildValueList{
										{Name: "name", Value: &ast.Value{Kind: ast.StringValue, Raw: "lastname"}},
									},
								}},
							},
						}},
					},
				},
			},
		},
	}

	node, err := ParseDirective(d, cat, "users", cols)
	if err != nil {
		t.Fatal(err)
	}
	if node.Op != "concat" || len(node.Args) != 3 {
		t.Fatalf("unexpected node: %+v", node)
	}
}
