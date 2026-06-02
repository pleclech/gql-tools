package gqlexpr

import (
	"encoding/json"
	"fmt"
)

// Node is a canonical expression AST node (JSON-stable).
type Node struct {
	Col   string `json:"col,omitempty"`
	Str   string `json:"str,omitempty"`
	Op    string `json:"op,omitempty"`
	Table string `json:"table,omitempty"`
	Args  []Node `json:"args,omitempty"`
}

// MarshalJSON stores expr nodes for NamedObject.ex.
func (n Node) MarshalJSON() ([]byte, error) {
	if n.Op != "" {
		return json.Marshal(struct {
			Op    string `json:"op"`
			Table string `json:"table,omitempty"`
			Args  []Node `json:"args"`
		}{Op: n.Op, Table: n.Table, Args: n.Args})
	}
	if n.Col != "" {
		return json.Marshal(struct {
			Col string `json:"col"`
		}{Col: n.Col})
	}
	if n.Str != "" {
		return json.Marshal(struct {
			Str string `json:"str"`
		}{Str: n.Str})
	}
	return nil, fmt.Errorf("gqlexpr: empty node")
}

// CollectColumns returns column names referenced in the tree.
func (n Node) CollectColumns() []string {
	var out []string
	n.collectColumns(&out)
	return out
}

func (n Node) collectColumns(out *[]string) {
	if n.Col != "" {
		*out = append(*out, n.Col)
	}
	for i := range n.Args {
		n.Args[i].collectColumns(out)
	}
}

// ParseNodeJSON decodes a stored expression payload.
func ParseNodeJSON(raw json.RawMessage) (Node, error) {
	var n Node
	if len(raw) == 0 {
		return n, fmt.Errorf("gqlexpr: empty expression")
	}
	if err := json.Unmarshal(raw, &n); err != nil {
		return n, err
	}
	if n.Op == "" && n.Col == "" && n.Str == "" {
		return n, fmt.Errorf("gqlexpr: invalid expression json")
	}
	return n, nil
}
