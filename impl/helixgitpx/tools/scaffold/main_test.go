package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRender_CreatesExpectedFiles(t *testing.T) {
	dst := t.TempDir()
	cfg := Config{
		Name:         "greet",
		ProtoPackage: "helixgitpx.greet.v1",
		HTTPPort:     8002,
		GRPCPort:     9002,
		HealthPort:   8082,
		Out:          dst,
	}
	if err := Render(cfg); err != nil {
		t.Fatalf("Render: %v", err)
	}

	expected := []string{
		"cmd/greet/main.go",
		"internal/app/app.go",
		"Makefile",
		"README.md",
		"go.mod",
		"deploy/Dockerfile",
	}
	for _, p := range expected {
		if _, err := os.Stat(filepath.Join(dst, p)); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}

	main, err := os.ReadFile(filepath.Join(dst, "cmd/greet/main.go"))
	if err != nil {
		t.Fatalf("read main: %v", err)
	}
	if !contains(string(main), "greet") {
		t.Errorf("main.go did not substitute name")
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
