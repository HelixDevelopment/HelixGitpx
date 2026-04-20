// Package log provides a structured, JSON logger built on log/slog.
//
// Options controls level, output, and service/version tags. FromContext/WithContext
// pass a child logger through a context.Context for request-scoped fields.
package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"
)

// Options configures New.
type Options struct {
	Level   string    // "debug"/"info"/"warn"/"error"
	Output  io.Writer // defaults to os.Stdout
	Service string    // added as "service" field on every record
	Version string    // added as "version" field on every record
}

// Logger wraps *slog.Logger with With/FromContext/WithContext helpers.
type Logger struct {
	sl *slog.Logger
}

type ctxKey struct{}

var global atomic.Pointer[Logger]

// New constructs a Logger.
func New(opts Options) *Logger {
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	handler := slog.NewJSONHandler(opts.Output, &slog.HandlerOptions{
		Level: parseLevel(opts.Level),
	})
	sl := slog.New(handler)
	if opts.Service != "" {
		sl = sl.With("service", opts.Service)
	}
	if opts.Version != "" {
		sl = sl.With("version", opts.Version)
	}
	lg := &Logger{sl: sl}
	global.Store(lg)
	return lg
}

// Default returns the most recently constructed Logger, or a noop logger if none.
func Default() *Logger {
	if lg := global.Load(); lg != nil {
		return lg
	}
	return &Logger{sl: slog.New(slog.NewJSONHandler(io.Discard, nil))}
}

// Info logs at info level with key/value pairs.
func (l *Logger) Info(msg string, kv ...any)  { l.sl.Info(msg, kv...) }
func (l *Logger) Warn(msg string, kv ...any)  { l.sl.Warn(msg, kv...) }
func (l *Logger) Error(msg string, kv ...any) { l.sl.Error(msg, kv...) }
func (l *Logger) Debug(msg string, kv ...any) { l.sl.Debug(msg, kv...) }

// With returns a child logger with the supplied fields attached.
func (l *Logger) With(kv ...any) *Logger { return &Logger{sl: l.sl.With(kv...)} }

// WithContext stores lg on ctx. Callers retrieve it with FromContext.
func WithContext(ctx context.Context, lg *Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, lg)
}

// FromContext returns the logger attached via WithContext, or Default().
func FromContext(ctx context.Context) *Logger {
	if lg, ok := ctx.Value(ctxKey{}).(*Logger); ok {
		return lg
	}
	return Default()
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
