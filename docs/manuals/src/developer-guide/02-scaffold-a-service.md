## 2. Scaffolding a new service

HelixGitpx's `tools/scaffold` generates a complete service skeleton
that matches the monorepo's conventions.

### 2.1 Generate

```bash
cd impl/helixgitpx
go run ./tools/scaffold service --name foobar
```

What you get:

```
services/foobar/
  cmd/foobar/main.go
  internal/app/app.go
  internal/domain/          (empty; add your invariants here)
  deploy/Dockerfile
  deploy/helm/Chart.yaml
  deploy/helm/values.yaml
  deploy/helm/templates/*   (deployment + service + configmap + sm + ingress + migrate-job)
  go.mod                    (joined to go.work)
  README.md
  Makefile
```

### 2.2 Add to go.work

The scaffold appends an entry automatically. If you regenerate the workspace:

```bash
cd impl/helixgitpx
go work sync
```

### 2.3 Register in Argo

Create `impl/helixgitpx-platform/argocd/applications/foobar.yaml`
cloning a sibling service's file and pointing `source.path` at
`impl/helixgitpx/services/foobar/deploy/helm`. Pick a sync-wave:

- **10** for business services.
- **7** for infra-style services (AI/search).
- **5** or earlier for foundational services.

### 2.4 Add a proto contract (if public)

```
impl/helixgitpx/proto/helixgitpx/foobar/v1/foobar.proto
```

Then:

```bash
cd impl/helixgitpx/proto
buf lint
buf generate
```

Generated Go lands in `impl/helixgitpx/gen/go/helixgitpx/foobar/v1/`.

### 2.5 Wire it up

- Domain logic goes in `internal/domain/` with TDD tests.
- HTTP / gRPC handlers go in `internal/handler/`.
- Composition root (`internal/app/app.go`) wires handlers into a
  `http.ServeMux` or Connect-RPC registrar, boots OTel, registers
  health endpoints, and handles graceful shutdown. See `services/hello/`
  for the reference implementation.

### 2.6 CI acceptance checks

Before opening a PR, run:

```bash
bash scripts/verify-everything.sh
bash scripts/verify-helm-charts.sh
bash scripts/verify-argo-paths.sh
(cd impl/helixgitpx && go test github.com/helixgitpx/helixgitpx/services/foobar/...)
```

All four must be green.

---
