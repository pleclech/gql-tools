package gqldoc

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/vektah/gqlparser/v2/ast"
)

// JSONWriter serializes GraphQL AST values and arguments as JSON.
type JSONWriter struct {
	w *bytes.Buffer
}

// NewJSONWriter creates a JSONWriter.
func NewJSONWriter() *JSONWriter {
	return &JSONWriter{w: &bytes.Buffer{}}
}

// Bytes returns the written JSON bytes.
func (j *JSONWriter) Bytes() []byte {
	return j.w.Bytes()
}

// String returns the written JSON string.
func (j *JSONWriter) String() string {
	return j.w.String()
}

func (j *JSONWriter) WriteChildValue(child *ast.ChildValue) error {
	if child == nil {
		return nil
	}
	w := j.w
	if child.Name != "" {
		w.WriteRune('"')
		w.WriteString(child.Name)
		w.WriteRune('"')
		w.WriteRune(':')
	}
	return j.WriteValue(child.Value)
}

func (j *JSONWriter) WriteValue(val *ast.Value) error {
	w := j.w
	if val == nil {
		w.WriteString("null")
		return nil
	}
	switch val.Kind {
	case ast.IntValue, ast.FloatValue, ast.BooleanValue:
		w.WriteString(val.Raw)
		return nil
	case ast.StringValue, ast.BlockValue, ast.EnumValue:
		quoted := strconv.Quote(val.Raw)
		w.WriteString(quoted)
		return nil
	case ast.NullValue:
		w.WriteString("null")
		return nil
	case ast.ListValue:
		w.WriteString("[")
		sep := ""
		for _, v := range val.Children {
			w.WriteString(sep)
			sep = ","
			if err := j.WriteChildValue(v); err != nil {
				return err
			}
		}
		w.WriteString("]")
		return nil
	case ast.ObjectValue:
		w.WriteString("{")
		sep := ""
		for _, v := range val.Children {
			w.WriteString(sep)
			sep = ","
			if err := j.WriteChildValue(v); err != nil {
				return err
			}
		}
		w.WriteString("}")
		return nil
	default:
		return fmt.Errorf("unknown value kind: %v", val.Kind)
	}
}

func (j *JSONWriter) WriteArgument(arg *ast.Argument) error {
	if arg == nil {
		return nil
	}
	w := j.w
	w.WriteRune('"')
	w.WriteString(arg.Name)
	w.WriteRune('"')
	w.WriteString(":")
	return j.WriteValue(arg.Value)
}

func (j *JSONWriter) WriteArgumentList(args ast.ArgumentList) error {
	if len(args) == 0 {
		return nil
	}
	w := j.w
	sep := ""
	w.WriteString("{")
	for _, arg := range args {
		w.WriteString(sep)
		sep = ","
		if err := j.WriteArgument(arg); err != nil {
			return err
		}
	}
	w.WriteString("}")
	return nil
}

func (j *JSONWriter) WriteDirective(directive *ast.Directive) error {
	if directive == nil {
		return nil
	}
	w := j.w
	w.WriteString(`{"@`)
	w.WriteString(directive.Name)
	w.WriteString(`":`)
	if err := j.WriteArgumentList(directive.Arguments); err != nil {
		return err
	}
	w.WriteString("}")
	return nil
}

func (j *JSONWriter) WriteDirectiveList(directives ast.DirectiveList) error {
	if len(directives) == 0 {
		return nil
	}
	w := j.w
	sep := ""
	w.WriteString("[")
	for _, directive := range directives {
		w.WriteString(sep)
		sep = ","
		if err := j.WriteDirective(directive); err != nil {
			return err
		}
	}
	w.WriteString("]")
	return nil
}

// Write serializes a supported AST expression as JSON.
func (j *JSONWriter) Write(expr interface{}) error {
	switch e := expr.(type) {
	case *ast.ChildValue:
		return j.WriteChildValue(e)
	case *ast.Value:
		return j.WriteValue(e)
	case *ast.Argument:
		return j.WriteArgument(e)
	case ast.ArgumentList:
		return j.WriteArgumentList(e)
	case *ast.Directive:
		return j.WriteDirective(e)
	case ast.DirectiveList:
		return j.WriteDirectiveList(e)
	default:
		return fmt.Errorf("unknown type %T", e)
	}
}
