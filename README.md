# gql-tools

Go libraries for working with GraphQL documents and related tooling.

## Modules

| Module | Description |
|---|---|
| [gqldoc](./gqldoc) | Parse, mutate, and format GraphQL query documents on [gqlparser/v2](https://github.com/vektah/gqlparser) |

## Install

Each package is a separate Go module:

```bash
go get github.com/pleclech/gql-tools/gqldoc@gqldoc/v0.1.0
```

## Development

```bash
go work init ./gqldoc
cd gqldoc && go test ./...
```

## License

MIT
