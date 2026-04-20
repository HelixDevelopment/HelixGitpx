// Package telemetry (pprof.go) exposes net/http/pprof handlers for continuous
// profiling via Pyroscope's pull-based scraping. RegisterPprof attaches the
// handlers to the passed mux (typically the health mux on a separate port).
package telemetry

import (
	"net/http"
	"net/http/pprof"
)

// RegisterPprof adds the standard pprof handlers to mux. Call it once per
// process from the composition root before serving the health mux.
func RegisterPprof(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}
