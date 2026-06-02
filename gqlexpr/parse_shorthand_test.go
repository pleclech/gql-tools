package gqlexpr

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestParseExprValueShorthand(t *testing.T) {
	cols := NewColumnRegistry(map[string][]string{
		"users_select_column": {"firstname", "lastname"},
	})

	colEnum, err := parseExprValue(&ast.Value{
		Kind: ast.ObjectValue,
		Children: ast.ChildValueList{
			{Name: "column", Value: &ast.Value{Kind: ast.EnumValue, Raw: "firstname"}},
		},
	}, "users", cols)
	if err != nil || colEnum.Col != "firstname" {
		t.Fatalf("column enum shorthand: %+v err=%v", colEnum, err)
	}

	colLegacy, err := parseExprValue(&ast.Value{
		Kind: ast.ObjectValue,
		Children: ast.ChildValueList{
			{Name: "column", Value: &ast.Value{
				Kind: ast.ObjectValue,
				Children: ast.ChildValueList{
					{Name: "name", Value: &ast.Value{Kind: ast.EnumValue, Raw: "lastname"}},
				},
			}},
		},
	}, "users", cols)
	if err != nil || colLegacy.Col != "lastname" {
		t.Fatalf("column legacy: %+v err=%v", colLegacy, err)
	}

	strShort, err := parseExprValue(&ast.Value{
		Kind: ast.ObjectValue,
		Children: ast.ChildValueList{
			{Name: "string", Value: &ast.Value{Kind: ast.StringValue, Raw: `" "`}},
		},
	}, "users", cols)
	if err != nil || strShort.Str != " " {
		t.Fatalf("string shorthand: %+v err=%v", strShort, err)
	}

	strLegacy, err := parseExprValue(&ast.Value{
		Kind: ast.ObjectValue,
		Children: ast.ChildValueList{
			{Name: "string", Value: &ast.Value{
				Kind: ast.ObjectValue,
				Children: ast.ChildValueList{
					{Name: "value", Value: &ast.Value{Kind: ast.StringValue, Raw: `"x"`}},
				},
			}},
		},
	}, "users", cols)
	if err != nil || strLegacy.Str != "x" {
		t.Fatalf("string legacy: %+v err=%v", strLegacy, err)
	}
}

func TestParseExprValueInvalidColumnEnum(t *testing.T) {
	cols := NewColumnRegistry(map[string][]string{
		"users_select_column": {"id"},
	})
	_, err := parseExprValue(&ast.Value{
		Kind: ast.ObjectValue,
		Children: ast.ChildValueList{
			{Name: "column", Value: &ast.Value{Kind: ast.EnumValue, Raw: "secret_col"}},
		},
	}, "users", cols)
	if err == nil {
		t.Fatal("expected disallowed column error")
	}
}
