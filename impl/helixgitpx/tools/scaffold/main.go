// Command scaffold renders a new HelixGitpx Go service from templates.
//
//	scaffold --name greet --proto helixgitpx.greet.v1 --http 8002 --grpc 9002
//
// When --dry-run is set, prints the list of files that would be written and exits 0.
package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed all:templates/service
var templatesFS embed.FS

// Config drives Render.
type Config struct {
	Name         string
	ProtoPackage string
	HTTPPort     int
	GRPCPort     int
	HealthPort   int
	Out          string
	DryRun       bool
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.Name, "name", "", "service name (lowercase, no spaces)")
	flag.StringVar(&cfg.ProtoPackage, "proto", "", "protobuf package (e.g. helixgitpx.greet.v1)")
	flag.IntVar(&cfg.HTTPPort, "http", 8000, "HTTP port")
	flag.IntVar(&cfg.GRPCPort, "grpc", 9000, "gRPC port")
	flag.IntVar(&cfg.HealthPort, "health", 8080, "health/metrics port")
	flag.StringVar(&cfg.Out, "out", "", "output directory (default: services/<name>)")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "list files without writing")
	flag.Parse()

	if cfg.Name == "" || cfg.ProtoPackage == "" {
		flag.Usage()
		os.Exit(2)
	}
	if cfg.Out == "" {
		cfg.Out = filepath.Join("services", cfg.Name)
	}
	if cfg.DryRun {
		if err := DryRun(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "scaffold: %v\n", err)
			os.Exit(1)
		}
		return
	}
	if err := Render(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "scaffold: %v\n", err)
		os.Exit(1)
	}
}

// Render writes all template files to cfg.Out, substituting values.
func Render(cfg Config) error {
	return walkTemplates(cfg, func(srcPath, dstPath string, raw []byte) error {
		if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
			return err
		}
		tpl, err := template.New(srcPath).Delims("<<", ">>").Parse(string(raw))
		if err != nil {
			return fmt.Errorf("parse %s: %w", srcPath, err)
		}
		f, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer f.Close()
		return tpl.Execute(f, cfg)
	})
}

// DryRun prints files that would be written.
func DryRun(cfg Config) error {
	return walkTemplates(cfg, func(_, dstPath string, _ []byte) error {
		fmt.Println(dstPath)
		return nil
	})
}

func walkTemplates(cfg Config, fn func(src, dst string, raw []byte) error) error {
	root := "templates/service"
	return fs.WalkDir(templatesFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		raw, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(path, root+"/")
		rel = strings.ReplaceAll(rel, "__name__", cfg.Name)
		rel = strings.TrimSuffix(rel, ".tmpl")
		return fn(path, filepath.Join(cfg.Out, rel), raw)
	})
}
