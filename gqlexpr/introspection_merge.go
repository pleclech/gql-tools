package gqlexpr

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// MergeIntrospectionResponse appends extension directives and types to a Hasura
// introspection JSON response (data.__schema.types / directives).
func MergeIntrospectionResponse(hasuraBody []byte, cat OperatorCatalog) ([]byte, error) {
	if cat == nil {
		return hasuraBody, nil
	}

	var root map[string]json.RawMessage
	if err := json.Unmarshal(hasuraBody, &root); err != nil {
		return nil, fmt.Errorf("introspection merge: %w", err)
	}
	dataRaw, ok := root["data"]
	if !ok {
		return hasuraBody, nil
	}

	var data map[string]json.RawMessage
	if err := json.Unmarshal(dataRaw, &data); err != nil {
		return nil, err
	}
	schemaRaw, ok := data["__schema"]
	if !ok {
		return hasuraBody, nil
	}

	var schema map[string]json.RawMessage
	if err := json.Unmarshal(schemaRaw, &schema); err != nil {
		return nil, err
	}

	enums := selectColumnEnumsFromIntrospectionSchema(schema)
	extraTypes := buildExtensionIntrospectionTypes(cat, enums)

	typesRaw, _ := schema["types"]
	var types []json.RawMessage
	if len(typesRaw) > 0 {
		_ = json.Unmarshal(typesRaw, &types)
	}
	existing := typeNamesSet(types)
	for _, t := range extraTypes {
		name, _ := typeNameFromIntrospectionType(t)
		if name != "" && existing[name] {
			continue
		}
		types = append(types, t)
		if name != "" {
			existing[name] = true
		}
	}

	typesJSON, err := json.Marshal(types)
	if err != nil {
		return nil, err
	}
	schema["types"] = typesJSON

	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	data["__schema"] = schemaJSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	root["data"] = dataJSON
	return json.Marshal(root)
}

func selectColumnEnumsFromIntrospectionSchema(schema map[string]json.RawMessage) map[string][]string {
	typesRaw := schema["types"]
	if len(typesRaw) == 0 {
		return nil
	}
	var types []map[string]json.RawMessage
	if err := json.Unmarshal(typesRaw, &types); err != nil {
		return nil
	}
	out := make(map[string][]string)
	for _, t := range types {
		kind := stringField(t, "kind")
		name := stringField(t, "name")
		if kind != "ENUM" || !strings.HasSuffix(name, "_select_column") {
			continue
		}
		var enumValues []struct {
			Name string `json:"name"`
		}
		if raw, ok := t["enumValues"]; ok {
			_ = json.Unmarshal(raw, &enumValues)
		}
		cols := make([]string, 0, len(enumValues))
		for _, ev := range enumValues {
			cols = append(cols, ev.Name)
		}
		out[name] = cols
	}
	return out
}

func buildExtensionIntrospectionDirectives(cat OperatorCatalog) []json.RawMessage {
	ops := catalogOps(cat)
	sort.Slice(ops, func(i, j int) bool { return ops[i].directive < ops[j].directive })
	out := make([]json.RawMessage, 0, len(ops))
	for _, op := range ops {
		argType := introspectionTypeRef("NON_NULL", introspectionTypeRef("SCALAR", nil, "ExtensionExpr", nil), "", nil)
		if op.list {
			argType = introspectionTypeRef("NON_NULL", introspectionTypeRef("LIST", introspectionTypeRef("NON_NULL", introspectionTypeRef("SCALAR", nil, "ExtensionExpr", nil), "", nil), "", nil), "", nil)
		}
		d := map[string]interface{}{
			"kind":        "DIRECTIVE",
			"name":        op.directive,
			"description": "Hasura extension expression operator (" + op.name + ")",
			"locations":   []string{"FIELD"},
			"args": []map[string]interface{}{
				{
					"name":         op.argName,
					"description":  nil,
					"type":         argType,
					"defaultValue": nil,
				},
			},
		}
		b, _ := json.Marshal(d)
		out = append(out, b)
	}
	return out
}

func buildExtensionIntrospectionTypes(cat OperatorCatalog, selectColumnEnums map[string][]string) []json.RawMessage {
	var out []json.RawMessage

	extensionExprScalar, _ := json.Marshal(map[string]interface{}{
		"kind":        "SCALAR",
		"name":        "ExtensionExpr",
		"description": "Opaque expression node; see per-table *Expr input types",
	})
	out = append(out, extensionExprScalar)

	ops := catalogOps(cat)
	tables := make([]string, 0, len(selectColumnEnums))
	for enumName := range selectColumnEnums {
		tables = append(tables, strings.TrimSuffix(enumName, "_select_column"))
	}
	sort.Strings(tables)

	for _, table := range tables {
		enumKey := table + "_select_column"
		cols := selectColumnEnums[enumKey]
		if len(cols) == 0 {
			continue
		}
		typeName := typeNameForTable(table)
		exprName := typeName + "Expr"
		columnEnum := enumKey // Hasura *_select_column enum

		for _, op := range ops {
			inName := typeName + op.titleName
			var fields []map[string]interface{}
			if op.list {
				fields = []map[string]interface{}{
					{"name": op.argName, "type": introspectionTypeRef("NON_NULL", introspectionTypeRef("LIST", introspectionTypeRef("NON_NULL", introspectionTypeRef("INPUT_OBJECT", nil, exprName, nil), "", nil), "", nil), "", nil), "defaultValue": nil},
				}
			} else {
				fields = []map[string]interface{}{
					{"name": op.argName, "type": introspectionTypeRef("NON_NULL", introspectionTypeRef("INPUT_OBJECT", nil, exprName, nil), "", nil), "defaultValue": nil},
				}
			}
			inType, _ := json.Marshal(map[string]interface{}{"kind": "INPUT_OBJECT", "name": inName, "inputFields": fields})
			out = append(out, inType)
		}

		exprFields := []map[string]interface{}{
			{"name": "column", "type": introspectionTypeRef("ENUM", nil, columnEnum, nil), "defaultValue": nil},
			{"name": "string", "type": introspectionTypeRef("SCALAR", nil, "String", nil), "defaultValue": nil},
		}
		for _, op := range ops {
			exprFields = append(exprFields, map[string]interface{}{
				"name": op.name, "type": introspectionTypeRef("INPUT_OBJECT", nil, typeName+op.titleName, nil), "defaultValue": nil,
			})
		}
		exprType, _ := json.Marshal(map[string]interface{}{"kind": "INPUT_OBJECT", "name": exprName, "inputFields": exprFields})
		out = append(out, exprType)
	}

	return out
}

func introspectionTypeRef(kind string, ofType interface{}, name string, _ interface{}) map[string]interface{} {
	m := map[string]interface{}{"kind": kind}
	if ofType != nil {
		m["ofType"] = ofType
	}
	if name != "" {
		m["name"] = name
	}
	return m
}

func typeNamesSet(types []json.RawMessage) map[string]bool {
	out := make(map[string]bool, len(types))
	for _, t := range types {
		if n, ok := typeNameFromIntrospectionType(t); ok {
			out[n] = true
		}
	}
	return out
}

func directiveNamesSet(directives []json.RawMessage) map[string]bool {
	out := make(map[string]bool, len(directives))
	for _, d := range directives {
		if n, ok := directiveNameFromIntrospection(d); ok {
			out[n] = true
		}
	}
	return out
}

func typeNameFromIntrospectionType(raw json.RawMessage) (string, bool) {
	var t struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &t); err != nil {
		return "", false
	}
	return t.Name, t.Name != ""
}

func directiveNameFromIntrospection(raw json.RawMessage) (string, bool) {
	var d struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(raw, &d); err != nil {
		return "", false
	}
	return d.Name, d.Name != ""
}

func stringField(m map[string]json.RawMessage, key string) string {
	var s string
	if raw, ok := m[key]; ok {
		_ = json.Unmarshal(raw, &s)
	}
	return s
}
