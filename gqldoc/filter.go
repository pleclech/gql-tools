package gqldoc

import (
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// FilterArguments returns arguments for which keep returns true.
func FilterArguments(args ast.ArgumentList, keep func(*ast.Argument) bool) ast.ArgumentList {
	if len(args) == 0 {
		return nil
	}
	out := make(ast.ArgumentList, 0, len(args))
	for _, arg := range args {
		if arg != nil && keep(arg) {
			out = append(out, arg)
		}
	}
	return out
}

// RemoveArgumentsByName removes arguments with the given names.
func RemoveArgumentsByName(args ast.ArgumentList, names ...string) ast.ArgumentList {
	remove := make(map[string]struct{}, len(names))
	for _, name := range names {
		remove[name] = struct{}{}
	}
	return FilterArguments(args, func(arg *ast.Argument) bool {
		_, drop := remove[arg.Name]
		return !drop
	})
}

// FilterDirectives returns directives for which keep returns true.
func FilterDirectives(dirs ast.DirectiveList, keep func(*ast.Directive) bool) ast.DirectiveList {
	if len(dirs) == 0 {
		return nil
	}
	out := make(ast.DirectiveList, 0, len(dirs))
	for _, d := range dirs {
		if d != nil && keep(d) {
			out = append(out, d)
		}
	}
	return out
}

// RemoveDirectivesByPrefix removes directives whose name starts with prefix.
func RemoveDirectivesByPrefix(dirs ast.DirectiveList, prefix string) ast.DirectiveList {
	return FilterDirectives(dirs, func(d *ast.Directive) bool {
		return !strings.HasPrefix(d.Name, prefix)
	})
}

// FilterVariableDefinitions returns variable definitions for which keep returns true.
func FilterVariableDefinitions(defs ast.VariableDefinitionList, keep func(*ast.VariableDefinition) bool) ast.VariableDefinitionList {
	if len(defs) == 0 {
		return nil
	}
	out := make(ast.VariableDefinitionList, 0, len(defs))
	for _, def := range defs {
		if def != nil && keep(def) {
			out = append(out, def)
		}
	}
	return out
}

// FilterSelectionSet returns selections for which keep returns true.
func FilterSelectionSet(ss ast.SelectionSet, keep func(ast.Selection) bool) ast.SelectionSet {
	if len(ss) == 0 {
		return nil
	}
	out := make(ast.SelectionSet, 0, len(ss))
	for _, sel := range ss {
		if keep(sel) {
			out = append(out, sel)
		}
	}
	return out
}
