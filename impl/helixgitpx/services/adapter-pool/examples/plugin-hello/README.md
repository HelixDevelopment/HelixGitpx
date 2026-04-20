# plugin-hello — example TinyGo WASM plugin

Demonstrates the `adapter-pool` WASM plugin ABI (ADR-0022).

## Build

```sh
tinygo build -o plugin-hello.wasm -target=wasi ./main.go
```

## ABI

The host calls exported WASM functions:
- `adapter_list_refs(src_ptr, src_len) -> (out_ptr, out_len)`
- `adapter_get_repo(src_ptr, src_len) -> (out_ptr, out_len)`
- `adapter_push(dst_ptr, dst_len, refs_ptr, refs_len) -> err_code`

All payloads are JSON-encoded `adapter.Source` / `adapter.RepoInfo` / etc. The host
provides host functions `helix_log(msg_ptr, msg_len)` and `helix_time_now() -> i64`.

M5 ships this example as a shape reference; runtime execution is enabled in M5 hardening.
