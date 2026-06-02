package gqlexpr

// NewCatalogAdapterFromOperators builds a CatalogAdapter from operator definitions.
func NewCatalogAdapterFromOperators(ops map[string]struct {
	Directive string
	SQL       string
	Args      []ArgSpec
}) *CatalogAdapter {
	by := make(map[string]struct {
		Name string
		SQL  string
		Args []ArgSpec
	}, len(ops))
	for name, op := range ops {
		by[op.Directive] = struct {
			Name string
			SQL  string
			Args []ArgSpec
		}{Name: name, SQL: op.SQL, Args: op.Args}
	}
	return &CatalogAdapter{ByDirective: by}
}

// SQLNameMap returns operator name -> SQL function name.
func SQLNameMap(adapter *CatalogAdapter) map[string]string {
	if adapter == nil {
		return nil
	}
	m := make(map[string]string, len(adapter.ByDirective))
	for _, op := range adapter.ByDirective {
		m[op.Name] = op.SQL
	}
	return m
}
