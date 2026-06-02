package gqlexpr

// IntrospectionType is a minimal Hasura introspection type entry.
type IntrospectionType struct {
	Kind       string
	Name       string
	EnumValues []struct {
		Name string
	}
}

// SelectColumnEnumsFromTypes extracts *_select_column enum values.
func SelectColumnEnumsFromTypes(types []IntrospectionType) map[string][]string {
	out := make(map[string][]string)
	for _, t := range types {
		if t.Kind != "ENUM" {
			continue
		}
		if len(t.EnumValues) == 0 {
			continue
		}
		if !hasSuffix(t.Name, "_select_column") {
			continue
		}
		cols := make([]string, 0, len(t.EnumValues))
		for _, ev := range t.EnumValues {
			cols = append(cols, ev.Name)
		}
		out[t.Name] = cols
	}
	return out
}

func hasSuffix(s, suf string) bool {
	return len(s) >= len(suf) && s[len(s)-len(suf):] == suf
}
