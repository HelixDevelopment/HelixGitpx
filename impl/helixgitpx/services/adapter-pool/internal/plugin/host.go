// Package plugin loads WASM plugins via wazero and dispatches adapter.Adapter
// calls to them through a stable ABI. ADR-0022 documents the wazero choice.
//
// The Host is a thin registry over compiled plugins. A concrete `Runtime`
// supplies the actual wazero-backed invocation; this package only manages
// plugin registration, caching, and lifecycle so the service layer can
// depend on a narrow interface.
package plugin

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrUnknownPlugin is returned when a plugin name has not been registered.
var ErrUnknownPlugin = errors.New("plugin: unknown name")

// ErrNoRuntime indicates that no runtime has been installed; the host cannot
// invoke anything until `SetRuntime` is called.
var ErrNoRuntime = errors.New("plugin: runtime not initialized")

// Runtime is the narrow surface an underlying WASM engine must satisfy.
// The production implementation under services/adapter-pool/internal/plugin/wazero/
// wraps tetratelabs/wazero; unit tests install a fake.
type Runtime interface {
	// Invoke calls `method` on the given compiled module with serialized args.
	Invoke(ctx context.Context, module CompiledModule, method string, args []byte) ([]byte, error)
	// Compile returns a CompiledModule ready to be invoked.
	Compile(ctx context.Context, wasm []byte) (CompiledModule, error)
	// Close releases any resources held by the runtime.
	Close(ctx context.Context) error
}

// CompiledModule is an opaque handle returned by the Runtime.Compile call.
type CompiledModule any

// Host manages registered plugins and dispatches calls to the installed runtime.
type Host struct {
	mu      sync.RWMutex
	runtime Runtime
	modules map[string]CompiledModule
}

// NewHost returns an empty Host. Install a runtime before Invoke.
func NewHost() *Host {
	return &Host{modules: map[string]CompiledModule{}}
}

// SetRuntime installs the wasm runtime used for subsequent compiles and
// invocations. Passing nil removes the runtime (used by tests for cleanup).
func (h *Host) SetRuntime(rt Runtime) {
	h.mu.Lock()
	h.runtime = rt
	h.mu.Unlock()
}

// Register compiles and installs a plugin under the given name. A plugin name
// must be globally unique within the host; re-registering the same name
// replaces the previous module.
func (h *Host) Register(ctx context.Context, name string, wasm []byte) error {
	h.mu.RLock()
	rt := h.runtime
	h.mu.RUnlock()
	if rt == nil {
		return ErrNoRuntime
	}
	mod, err := rt.Compile(ctx, wasm)
	if err != nil {
		return fmt.Errorf("plugin %s: compile: %w", name, err)
	}
	h.mu.Lock()
	h.modules[name] = mod
	h.mu.Unlock()
	return nil
}

// Invoke dispatches a named method on a plugin with marshalled args.
// The ABI is documented in examples/plugin-hello/README.md.
func (h *Host) Invoke(ctx context.Context, pluginName, method string, args []byte) ([]byte, error) {
	h.mu.RLock()
	rt := h.runtime
	mod, ok := h.modules[pluginName]
	h.mu.RUnlock()
	if rt == nil {
		return nil, ErrNoRuntime
	}
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownPlugin, pluginName)
	}
	return rt.Invoke(ctx, mod, method, args)
}

// Close releases runtime resources. Safe to call on a nil Host.
func (h *Host) Close(ctx context.Context) error {
	if h == nil {
		return nil
	}
	h.mu.Lock()
	rt := h.runtime
	h.modules = map[string]CompiledModule{}
	h.runtime = nil
	h.mu.Unlock()
	if rt == nil {
		return nil
	}
	return rt.Close(ctx)
}
