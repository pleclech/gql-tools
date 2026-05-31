# gqlhasura

Reusable Hasura GraphQL filter helpers built on [gqldoc](../gqldoc).

Part of [gql-tools](https://github.com/pleclech/gql-tools).

Project-specific transforms (e.g. `x_*` directive encoding) belong in the consuming app, not here.

## Install

```bash
go get github.com/pleclech/gql-tools/gqlhasura@gqlhasura/v0.1.0
```

## API overview

| Area | Functions |
|---|---|
| Where fixes | `NormalizeOrToList`, `NormalizeWhereOrLists` |

## Example

```go
import (
    "github.com/pleclech/gql-tools/gqlhasura"
    "github.com/pleclech/gql-tools/gqldoc"
)

gqldoc.MapWhereArguments(op, func(_ *ast.Field, where *ast.Value) error {
    gqlhasura.NormalizeWhereOrLists(where)
    return nil
})
```

## License

MIT
