package gqldoc

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestParseQuery(t *testing.T) {
	doc, err := ParseQuery(`query { hello }`)
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Operations) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(doc.Operations))
	}
}

func TestParseWhereLiteral(t *testing.T) {
	t.Run("bare fragment", func(t *testing.T) {
		node, err := ParseWhereLiteral(`status:{code:{_eq:"active"}}`)
		if err != nil {
			t.Fatal(err)
		}
		if node.Name != "status" {
			t.Fatalf("expected status, got %q", node.Name)
		}
		if node.Value == nil || node.Value.Kind != ast.ObjectValue {
			t.Fatalf("expected object value, got %+v", node.Value)
		}
	})

	t.Run("wrapped fragment", func(t *testing.T) {
		node, err := ParseWhereLiteral(`{status:{code:{_eq:"active"}}}`)
		if err != nil {
			t.Fatal(err)
		}
		if node.Name != "status" {
			t.Fatalf("expected status, got %q", node.Name)
		}
	})

	t.Run("empty fragment", func(t *testing.T) {
		if _, err := ParseWhereLiteral(""); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestParseDirectiveJSON(t *testing.T) {
	tests := []struct {
		name string
		qry  string
		want string
	}{
		{
			name: "array",
			qry:  `@x_having(_and:[{_count:{on:email,where:{_gt:1}}}, {_sum:{on:price,where:{_is_null:true}}}])`,
			want: `{"@x_having":{"_and":[{"_count":{"on":"email","where":{"_gt":1}}},{"_sum":{"on":"price","where":{"_is_null":true}}}]}}`,
		},
		{
			name: "simple",
			qry:  `@x_having(_count:{on:email,where:{_gt:1}})`,
			want: `{"@x_having":{"_count":{"on":"email","where":{"_gt":1}}}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDirective(tt.qry)
			if err != nil {
				t.Fatal(err)
			}
			jw := NewJSONWriter()
			if err := jw.WriteDirective(got); err != nil {
				t.Fatal(err)
			}
			if jw.String() != tt.want {
				t.Fatalf("got %q want %q", jw.String(), tt.want)
			}
		})
	}
}

func TestFormatOperation(t *testing.T) {
	doc, err := ParseQuery(`query Q { a b { c } }`)
	if err != nil {
		t.Fatal(err)
	}
	got := FormatOperation(doc.Operations[0])
	if got == "" {
		t.Fatal("expected formatted operation")
	}
}
