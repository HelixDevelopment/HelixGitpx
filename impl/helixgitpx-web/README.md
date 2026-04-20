# helixgitpx-web (Nx workspace, Angular 19)

## Quick start

```sh
npm install
npx nx serve web
```

The root `make gen` regenerates `libs/proto/` from `impl/helixgitpx/proto/`.

## Layout

- `apps/web/` — the Angular shell (M6 expands with real screens).
- `libs/proto/` — protobuf + Connect-ES codegen (committed).
- `libs/ui/` — design system (M6).
- `libs/data/` — data-access layers (M6).
