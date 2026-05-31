package gqldoc

import (
	"bytes"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

// Format prints a query document as GraphQL source.
func Format(doc *ast.QueryDocument) string {
	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatQueryDocument(doc)
	return buf.String()
}

// FormatOperation prints a single operation as a query document.
func FormatOperation(op *ast.OperationDefinition) string {
	return Format(&ast.QueryDocument{Operations: []*ast.OperationDefinition{op}})
}
