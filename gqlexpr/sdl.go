package gqlexpr

import (
	"fmt"
	"sort"
	"strings"
)

// BuildExtensionSDL generates GraphQL SDL documenting per-table expression inputs
// and extension directives (for IDE discovery / monaco-graphql secondary schema).
func BuildExtensionSDL(cat OperatorCatalog, selectColumnEnums map[string][]string) string {
	var b strings.Builder
	b.WriteString("# Hasura extension expression schema (generated)\n")
	b.WriteString("# Expr atoms: { column: COL } or { string: \"...\" } (legacy nested forms still accepted by the proxy).\n\n")

	ops := catalogOps(cat)
	sort.Slice(ops, func(i, j int) bool { return ops[i].directive < ops[j].directive })

	tables := make([]string, 0, len(selectColumnEnums))
	for enumName := range selectColumnEnums {
		tables = append(tables, strings.TrimSuffix(enumName, "_select_column"))
	}
	sort.Strings(tables)

	for _, table := range tables {
		enumKey := table + "_select_column"
		cols, ok := selectColumnEnums[enumKey]
		if !ok {
			cols = selectColumnEnums[table]
		}
		if len(cols) == 0 {
			continue
		}
		typeName := typeNameForTable(table)
		exprName := typeName + "Expr"
		columnEnum := enumKey

		for _, op := range ops {
			inName := typeName + op.titleName
			if op.list {
				b.WriteString(fmt.Sprintf("input %s {\n  %s: [%s!]!\n}\n\n", inName, op.argName, exprName))
			} else {
				b.WriteString(fmt.Sprintf("input %s {\n  %s: %s!\n}\n\n", inName, op.argName, exprName))
			}
		}

		b.WriteString(fmt.Sprintf("input %s {\n", exprName))
		b.WriteString(fmt.Sprintf("  column: %s\n", columnEnum))
		b.WriteString("  string: String\n")
		for _, op := range ops {
			b.WriteString(fmt.Sprintf("  %s: %s\n", op.name, typeName+op.titleName))
		}
		b.WriteString("}\n\n")

		// Table-scoped directive documentation (directives share names; args reference this table's Expr list).
		for _, op := range ops {
			if op.list {
				b.WriteString(fmt.Sprintf(
					"# On type %s fields: @%s(%s: [{ column: COL } | { string: \"...\" }])\n",
					table, op.directive, op.argName,
				))
			}
		}
	}

	// Global directive declarations (argument type is documented per table above).
	for _, op := range ops {
		if op.list {
			b.WriteString(fmt.Sprintf(
				"directive @%s(%s: [%s!]!) repeatable on FIELD\n\n",
				op.directive, op.argName, "ExtensionExpr",
			))
		} else {
			b.WriteString(fmt.Sprintf(
				"directive @%s(%s: %s!) repeatable on FIELD\n\n",
				op.directive, op.argName, "ExtensionExpr",
			))
		}
	}

	b.WriteString("scalar ExtensionExpr\n")

	return b.String()
}

// CatalogOpInfo describes one expression operator for introspection and parsing.
type CatalogOpInfo struct {
	Directive string
	Name      string
	TitleName string
	ArgName   string
	List      bool
}

type catalogOp struct {
	directive string
	name      string
	titleName string
	argName   string
	list      bool
}

func catalogOps(cat OperatorCatalog) []catalogOp {
	ca, ok := cat.(*CatalogAdapter)
	if !ok || ca == nil {
		return nil
	}
	out := make([]catalogOp, 0, len(ca.ByDirective))
	for directive, op := range ca.ByDirective {
		e := catalogOp{
			directive: directive,
			name:      op.Name,
			titleName: title(op.Name),
		}
		if len(op.Args) > 0 {
			e.argName = op.Args[0].Name
			e.list = op.Args[0].Kind == "expr_list"
		}
		out = append(out, e)
	}
	return out
}

// CatalogOperatorInfos returns operators from a catalog (sorted by directive name).
func CatalogOperatorInfos(cat OperatorCatalog) []CatalogOpInfo {
	ops := catalogOps(cat)
	sort.Slice(ops, func(i, j int) bool { return ops[i].directive < ops[j].directive })
	out := make([]CatalogOpInfo, len(ops))
	for i, op := range ops {
		out[i] = CatalogOpInfo{
			Directive: op.directive,
			Name:      op.name,
			TitleName: op.titleName,
			ArgName:   op.argName,
			List:      op.list,
		}
	}
	return out
}

// TableGraphQLName maps a Hasura table name (e.g. users) to a GraphQL type prefix (Users).
func TableGraphQLName(table string) string {
	return typeNameForTable(table)
}

func typeNameForTable(table string) string {
	parts := strings.Split(table, "_")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

func title(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
