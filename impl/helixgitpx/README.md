# helixgitpx (Go monorepo)

Contains the Go implementation of HelixGitpx: shared `platform` libraries, service binaries, the `scaffold` tool, and the proto root.

## Layout

```
├── go.work                Go 1.23 workspace
├── platform/              shared libraries (14 packages)
├── services/              service binaries (hello in M1; more in M3+)
├── tools/scaffold/        service-template renderer (Go binary, no Python)
├── proto/                 protobuf sources (buf module buf.build/helixgitpx/core)
├── gen/go/                generated Go code (committed; see .gitattributes)
└── api/openapi/           generated OpenAPI (committed)
```

## Daily commands

```sh
make lint     # golangci-lint + buf lint
make test     # go test -race
make gen      # regenerate protobuf bindings
make fmt      # gofumpt + goimports
```

See the shared lib docs under `platform/*/README.md`.
