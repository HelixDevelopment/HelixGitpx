//go:build tinygo_wasi

// Package main is an example TinyGo WASM plugin for the HelixGitpx adapter-pool.
// Build: tinygo build -o plugin-hello.wasm -target=wasi ./main.go
package main

//export adapter_get_repo
func adapter_get_repo(_ uintptr, _ uint32) (uint32, uint32) {
	// Return a fixed JSON payload: {"default":"main","private":false}.
	// Real pointer/length return omitted; this is a shape-only example.
	return 0, 0
}

func main() {}
