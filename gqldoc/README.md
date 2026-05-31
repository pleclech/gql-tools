# gqldoc

Helpers to parse, mutate, and format GraphQL query documents on top of [gqlparser/v2](https://github.com/vektah/gqlparser).

Part of [gql-tools](https://github.com/pleclech/gql-tools).

## Install

```bash
go get github.com/pleclech/gql-tools/gqldoc@gqldoc/v0.1.0
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/pleclech/gql-tools/gqldoc"
	"github.com/vektah/gqlparser/v2/ast"
)

func main() {
	doc, err := gqldoc.ParseQuery(`query { users(where: { id: { _eq: 1 } }) { id name } }`)
	if err != nil {
		panic(err)
	}

	_ = gqldoc.MapWhereArguments(doc.Operations[0], func(_ *ast.Field, where *ast.Value) error {
		gqldoc.MapObjectFields(where, func(path gqldoc.ObjectPath) (*ast.ChildValue, bool) {
			return path.Node, false
		})
		return nil
	})

	fmt.Println(gqldoc.Format(doc))
}
```

## API overview

| Area | Functions |
|---|---|
| Parse / format | `ParseQuery`, `ParseWrapped`, `ParseWhereLiteral`, `ParseValueLiteral`, `Format`, `FormatOperation` |
| Navigation | `AsField`, `FirstOperation`, `FirstField`, `FieldResponseName`, `FindArgument`, `FindWhere` |
| Builders | `Obj`, `Field`, `Str`, `Int`, `Bool`, `Null`, `Var`, `List`, `WrapNot` |
| Filters | `FilterArguments`, `RemoveDirectivesByPrefix`, `FilterVariableDefinitions` |
| Walking | `WalkFields`, `MapWhereArguments` |
| Where toolkit | `WalkObjectFields`, `CollectVariables`, `CollectVariablePaths`, `ReplaceAt`, `ApplyOptionalVariableFallback` |
| Comments | `CommentText`, `ParseCommentPrefix`, `ChildComment` |
| JSON | `JSONWriter` |

## Local development with a consumer

When working on gqldoc alongside a project that depends on it, add temporarily to the consumer `go.mod`:

```go
replace github.com/pleclech/gql-tools/gqldoc => ../gql-tools/gqldoc
```

Do not commit the replace directive.

## License

MIT
