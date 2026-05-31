// Package gqldoc provides helpers to parse, mutate, and format GraphQL query documents
// on top of github.com/vektah/gqlparser/v2.
//
// Typical usage:
//
//	doc, err := gqldoc.ParseQuery(`query { users(where: { id: { _eq: 1 } }) { id } }`)
//	gqldoc.MapWhereArguments(doc.Operations[0], func(field *ast.Field, where *ast.Value) error {
//	    gqldoc.MapObjectFields(where, func(path gqldoc.ObjectPath) (*ast.ChildValue, bool) {
//	        return path.Node, false
//	    })
//	    return nil
//	})
//	fmt.Println(gqldoc.Format(doc))
package gqldoc
