package gqldoc

import (
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// CommentText merges all lines in a comment group into a single string.
func CommentText(cg *ast.CommentGroup) string {
	if cg == nil || len(cg.List) == 0 {
		return ""
	}
	var b strings.Builder
	for _, comment := range cg.List {
		b.WriteString(strings.TrimSpace(comment.Text()))
		b.WriteString("\n")
	}
	return b.String()
}

// ParseCommentPrefix returns the payload after prefix when the merged comment starts with it.
func ParseCommentPrefix(cg *ast.CommentGroup, prefix string) (string, bool) {
	text := strings.TrimSpace(CommentText(cg))
	if text == "" {
		return "", false
	}
	after, ok := strings.CutPrefix(text, prefix)
	if !ok {
		return "", false
	}
	return strings.TrimSpace(after), true
}

// ChildComment returns the comment on a child value or its nested value.
func ChildComment(cv *ast.ChildValue) *ast.CommentGroup {
	if cv == nil {
		return nil
	}
	if cv.Comment != nil {
		return cv.Comment
	}
	if cv.Value != nil {
		return cv.Value.Comment
	}
	return nil
}
