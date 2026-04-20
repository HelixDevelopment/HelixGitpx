// Package plugin loads WASM plugins via wazero and dispatches adapter.Adapter
// calls to them through a stable ABI. ADR-0022 documents the wazero choice.
//
// M5 ships the host shell; full ABI wiring arrives when the first real plugin
// lands (the example under ../../examples/plugin-hello/ demonstrates the boundary).
package plugin

import (
	"context"
	"errors"
)

// Host is a WASM plugin runtime. Real implementation uses tetratelabs/wazero
// to compile + instantiate .wasm modules fetched from Vault KV at kv/plugins/<name>.
type Host struct{}

// NewHost returns a Host. Loading plugins is deferred to M5 hardening.
func NewHost() *Host { return &Host{} }

// Invoke dispatches a named method on a plugin with marshalled args. The ABI
// is documented in examples/plugin-hello/README.md.
func (h *Host) Invoke(_ context.Context, pluginName, method string, args []byte) ([]byte, error) {
	// TODO(M5-hardening): wire wazero runtime + plugin lookup.
	_ = pluginName
	_ = method
	_ = args
	return nil, errors.New("plugin: wazero runtime not initialized (M5 hardening deferred)")
}
