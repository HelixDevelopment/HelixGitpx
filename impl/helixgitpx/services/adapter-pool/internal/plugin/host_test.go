package plugin

import (
	"context"
	"errors"
	"testing"
)

type fakeRuntime struct {
	compiled [][]byte
	invoked  []string
	closed   bool
}

func (f *fakeRuntime) Compile(_ context.Context, wasm []byte) (CompiledModule, error) {
	f.compiled = append(f.compiled, wasm)
	return string(wasm), nil
}

func (f *fakeRuntime) Invoke(_ context.Context, mod CompiledModule, method string, args []byte) ([]byte, error) {
	f.invoked = append(f.invoked, method)
	return []byte(mod.(string) + ":" + method + ":" + string(args)), nil
}

func (f *fakeRuntime) Close(_ context.Context) error {
	f.closed = true
	return nil
}

func TestHost_InvokeWithoutRuntime(t *testing.T) {
	h := NewHost()
	_, err := h.Invoke(context.Background(), "any", "m", nil)
	if !errors.Is(err, ErrNoRuntime) {
		t.Fatalf("want ErrNoRuntime, got %v", err)
	}
}

func TestHost_UnknownPlugin(t *testing.T) {
	h := NewHost()
	h.SetRuntime(&fakeRuntime{})
	_, err := h.Invoke(context.Background(), "unknown", "m", nil)
	if !errors.Is(err, ErrUnknownPlugin) {
		t.Fatalf("want ErrUnknownPlugin, got %v", err)
	}
}

func TestHost_RegisterAndInvoke(t *testing.T) {
	fake := &fakeRuntime{}
	h := NewHost()
	h.SetRuntime(fake)

	if err := h.Register(context.Background(), "hello", []byte("wasm-bytes")); err != nil {
		t.Fatalf("register: %v", err)
	}
	out, err := h.Invoke(context.Background(), "hello", "greet", []byte("world"))
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if got, want := string(out), "wasm-bytes:greet:world"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}

	if err := h.Close(context.Background()); err != nil {
		t.Fatalf("close: %v", err)
	}
	if !fake.closed {
		t.Fatal("runtime Close not invoked")
	}
}
