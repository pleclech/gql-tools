package gqldoc

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)

// ObjectPath identifies a child field inside an object value tree.
type ObjectPath struct {
	Root   *ast.Value
	Parent *ast.Value
	Index  int
	Node   *ast.ChildValue
}

// WalkObjectFields traverses all named fields in an object value tree depth-first.
func WalkObjectFields(root *ast.Value, fn func(ObjectPath) error) error {
	if root == nil || root.Kind != ast.ObjectValue {
		return nil
	}
	var walk func(parent *ast.Value, index int, node *ast.ChildValue) error
	walk = func(parent *ast.Value, index int, node *ast.ChildValue) error {
		if node == nil {
			return nil
		}
		path := ObjectPath{
			Root:   root,
			Parent: parent,
			Index:  index,
			Node:   node,
		}
		if err := fn(path); err != nil {
			return err
		}
		if node.Value == nil {
			return nil
		}
		switch node.Value.Kind {
		case ast.ObjectValue:
			for i, child := range node.Value.Children {
				if err := walk(node.Value, i, child); err != nil {
					return err
				}
			}
		case ast.ListValue:
			for i, child := range node.Value.Children {
				if err := walk(node.Value, i, child); err != nil {
					return err
				}
			}
		}
		return nil
	}
	for i, child := range root.Children {
		if err := walk(nil, i, child); err != nil {
			return err
		}
	}
	return nil
}

// MapObjectFields replaces object fields when fn returns ok=true.
func MapObjectFields(root *ast.Value, fn func(ObjectPath) (*ast.ChildValue, bool)) {
	if root == nil || root.Kind != ast.ObjectValue {
		return
	}
	_ = WalkObjectFields(root, func(path ObjectPath) error {
		if replacement, ok := fn(path); ok {
			ReplaceAt(path, replacement)
		}
		return nil
	})
}

// ReplaceAt replaces the child value at path with replacement.
func ReplaceAt(path ObjectPath, replacement *ast.ChildValue) {
	if path.Node == nil || replacement == nil {
		return
	}
	if path.Parent == nil {
		path.Root.Children[path.Index] = replacement
		return
	}
	path.Parent.Children[path.Index] = replacement
}

// CollectVariables returns all variable references in an object value tree keyed by variable name.
func CollectVariables(root *ast.Value) map[string][]ObjectPath {
	out := make(map[string][]ObjectPath)
	if root == nil {
		return out
	}
	_ = WalkObjectFields(root, func(path ObjectPath) error {
		if path.Node == nil || path.Node.Value == nil {
			return nil
		}
		value := path.Node.Value
		if value.Kind == ast.Variable && value.Raw != "" {
			name := VariableName(value.Raw)
			pathCopy := path
			out[name] = append(out[name], pathCopy)
		}
		return nil
	})
	return out
}

// VariablePathEntry is one segment on the path from a where root to a variable reference.
type VariablePathEntry struct {
	ObjectPath
	Fallback *ast.ChildValue
}

// VariablePath is a full path to a variable use, including per-segment fallbacks.
type VariablePath struct {
	Segments []VariablePathEntry
}

// CollectVariablePaths records variable uses with segment chains and optional fallbacks.
func CollectVariablePaths(root *ast.Value, fallbackFor func(ObjectPath) (*ast.ChildValue, error)) (map[string][]VariablePath, error) {
	out := make(map[string][]VariablePath)
	if root == nil || root.Kind != ast.ObjectValue {
		return out, nil
	}

	var walk func(segments []VariablePathEntry, parent *ast.Value, index int, node *ast.ChildValue) error
	walk = func(segments []VariablePathEntry, parent *ast.Value, index int, node *ast.ChildValue) error {
		if node == nil {
			return nil
		}
		path := ObjectPath{
			Root:   root,
			Parent: parent,
			Index:  index,
			Node:   node,
		}
		entry := VariablePathEntry{ObjectPath: path}
		if fallbackFor != nil {
			fallback, err := fallbackFor(path)
			if err != nil {
				return err
			}
			entry.Fallback = fallback
		}
		segments = append(segments, entry)

		if node.Value != nil && node.Value.Kind == ast.Variable && node.Value.Raw != "" {
			name := VariableName(node.Value.Raw)
			copySegments := append([]VariablePathEntry(nil), segments...)
			out[name] = append(out[name], VariablePath{Segments: copySegments})
		}

		if node.Value != nil {
			for i, child := range node.Value.Children {
				if err := walk(segments, node.Value, i, child); err != nil {
					return err
				}
			}
		}
		return nil
	}

	for i, child := range root.Children {
		if err := walk(nil, nil, i, child); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// ApplyOptionalVariableFallback replaces the appropriate ancestor of a variable path with its fallback.
func ApplyOptionalVariableFallback(entry VariablePath) error {
	segments := entry.Segments
	if len(segments) == 0 {
		return nil
	}

	i := len(segments) - 1
	for ; i >= 0; i-- {
		value := segments[i].Node.Value
		if value == nil {
			break
		}
		kind := value.Kind
		if kind != ast.ObjectValue && kind != ast.Variable {
			break
		}
	}
	i++
	if i >= len(segments) {
		i = len(segments) - 1
	}
	if i < 0 {
		return nil
	}

	parent := segments[i]
	parentValue := parent.Node.Value
	if parentValue != nil && parentValue.Kind == ast.ListValue {
		return fmt.Errorf("not implemented (list value)")
	}
	if parent.Parent != nil && parent.Parent.Kind == ast.ListValue {
		return fmt.Errorf("not implemented (list value)")
	}

	fallback := parent.Fallback
	if fallback == nil {
		fallback = Field(parent.Node.Name, Obj())
	}
	ReplaceAt(parent.ObjectPath, fallback)
	return nil
}
