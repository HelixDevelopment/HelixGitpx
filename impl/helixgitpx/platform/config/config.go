// Package config loads typed configuration structs from environment variables
// with struct-tag-driven defaults, required markers, and list splitting.
//
// Struct tags supported:
//
//	env:"NAME"        environment variable name (required if field not anonymous)
//	default:"value"   default when env var unset
//	required:"true"   return error when unset and no default
//	split:","         for []string, split on the given separator
//
// Precedence: env var > default. File-based config and CLI flags may be layered
// by callers; this package intentionally keeps the minimal surface.
package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Options control Load.
type Options struct {
	// Prefix is prepended to every env tag as "<Prefix>_<tag>" when set.
	Prefix string
}

// Load populates *dst from the environment.
func Load(dst any, opts Options) error {
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config: Load requires a pointer to a struct, got %T", dst)
	}
	return loadStruct(v.Elem(), opts.Prefix)
}

func loadStruct(v reflect.Value, prefix string) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		fv := v.Field(i)

		if sf.Type.Kind() == reflect.Struct && sf.Anonymous {
			if err := loadStruct(fv, prefix); err != nil {
				return err
			}
			continue
		}

		envName := sf.Tag.Get("env")
		if envName == "" {
			continue
		}
		if prefix != "" {
			envName = prefix + "_" + envName
		}

		raw, present := os.LookupEnv(envName)
		if !present {
			raw = sf.Tag.Get("default")
		}
		if raw == "" && sf.Tag.Get("required") == "true" {
			return fmt.Errorf("config: required env %s is unset", envName)
		}
		if raw == "" {
			continue
		}
		if err := assign(fv, raw, sf.Tag.Get("split")); err != nil {
			return fmt.Errorf("config: %s: %w", envName, err)
		}
	}
	return nil
}

func assign(fv reflect.Value, raw, split string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(raw)
	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}
		fv.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fv.Type() == reflect.TypeOf(time.Duration(0)) {
			d, err := time.ParseDuration(raw)
			if err != nil {
				return err
			}
			fv.SetInt(int64(d))
			return nil
		}
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return err
		}
		fv.SetUint(n)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(f)
	case reflect.Slice:
		if split == "" {
			split = ","
		}
		parts := strings.Split(raw, split)
		out := reflect.MakeSlice(fv.Type(), len(parts), len(parts))
		for i, p := range parts {
			if err := assign(out.Index(i), strings.TrimSpace(p), ""); err != nil {
				return err
			}
		}
		fv.Set(out)
	default:
		return fmt.Errorf("unsupported kind %s", fv.Kind())
	}
	return nil
}
