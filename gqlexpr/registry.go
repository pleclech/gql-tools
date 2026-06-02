package gqlexpr

import "strings"

// ColumnRegistry maps GraphQL table/type names to allowed column names
// (from Hasura *_select_column enums).
type ColumnRegistry struct {
	// ColumnsByTable keys: type name without _select_column suffix (e.g. barnes_users).
	ColumnsByTable map[string][]string
}

// NewColumnRegistry builds a registry from introspection enum type names and values.
func NewColumnRegistry(selectColumnEnums map[string][]string) *ColumnRegistry {
	r := &ColumnRegistry{ColumnsByTable: make(map[string][]string, len(selectColumnEnums))}
	for enumName, cols := range selectColumnEnums {
		table := strings.TrimSuffix(enumName, "_select_column")
		r.ColumnsByTable[table] = cols
	}
	return r
}

// Allowed reports whether col is valid for tableContext.
func (r *ColumnRegistry) Allowed(tableContext, col string) bool {
	if r == nil {
		return false
	}
	cols, ok := r.ColumnsByTable[tableContext]
	if !ok {
		return false
	}
	for _, c := range cols {
		if c == col {
			return true
		}
	}
	return false
}

// ResolveTableContext maps a GraphQL object type or root field name to a registry key.
func (r *ColumnRegistry) ResolveTableContext(name string) string {
	if r == nil {
		return name
	}
	if _, ok := r.ColumnsByTable[name]; ok {
		return name
	}
	// Hasura often prefixes schema: barnes_users vs users
	for key := range r.ColumnsByTable {
		if strings.HasSuffix(key, "_"+name) || key == name {
			return key
		}
	}
	return name
}
