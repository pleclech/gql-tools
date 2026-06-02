package gqlexpr

import (
	"strings"
	"testing"
)

func TestMergeIntrospectionResponse_addsExpressionDirective(t *testing.T) {
	hasura := []byte(`{
		"data": {
			"__schema": {
				"types": [
					{
						"kind": "ENUM",
						"name": "users_select_column",
						"enumValues": [
							{"name": "id"},
							{"name": "firstname"}
						]
					}
				],
				"directives": []
			}
		}
	}`)

	cat := NewCatalogAdapterFromOperators(map[string]struct {
		Directive string
		SQL       string
		Args      []ArgSpec
	}{
		"concat": {
			Directive: "x_concat",
			SQL:       "concat",
			Args:      []ArgSpec{{Name: "parts", Kind: "expr_list"}},
		},
	})

	merged, err := MergeIntrospectionResponse(hasura, cat)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(merged), `"name":"UsersExpr"`) {
		t.Fatalf("expected UsersExpr type in merged schema")
	}
	if !strings.Contains(string(merged), `"name":"UsersConcat"`) {
		t.Fatalf("expected UsersConcat input type in merged schema")
	}
}
