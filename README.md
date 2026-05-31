# gql-tools

Go libraries for working with GraphQL documents and related tooling.

## Modules

| Module | Description |
|---|---|
| [gqldoc](./gqldoc) | Parse, mutate, and format GraphQL query documents on [gqlparser/v2](https://github.com/vektah/gqlparser) |
| [gqlhasura](./gqlhasura) | Reusable Hasura filter helpers (`_or` list normalization) |

## Install

Each package is a separate Go module:

```bash
go get github.com/pleclech/gql-tools/gqldoc@gqldoc/v0.1.0
go get github.com/pleclech/gql-tools/gqlhasura@gqlhasura/v0.1.0
```

## Development

```bash
go work init ./gqldoc ./gqlhasura
cd gqlhasura && go test ./...
```

## License

MIT
