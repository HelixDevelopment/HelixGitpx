// Package gin builds a Gin router preconfigured with HelixGitpx middleware:
// recovery, request-id, service/version headers, structured logging.
package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/helixgitpx/platform/log"
)

// Options configures NewRouter.
type Options struct {
	Service string
	Version string
	Mode    string // "debug", "release", "test"; defaults to "release"
}

// NewRouter constructs an *gin.Engine.
func NewRouter(opts Options) *gin.Engine {
	if opts.Mode == "" {
		opts.Mode = gin.ReleaseMode
	}
	gin.SetMode(opts.Mode)
	r := gin.New()
	r.Use(
		gin.Recovery(),
		identityHeaders(opts),
		loggingMiddleware(),
	)
	return r
}

func identityHeaders(opts Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		if opts.Service != "" {
			c.Header("X-HelixGitpx-Service", opts.Service)
		}
		if opts.Version != "" {
			c.Header("X-HelixGitpx-Version", opts.Version)
		}
		c.Next()
	}
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		lg := log.FromContext(c.Request.Context()).With(
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
		)
		if c.Writer.Status() >= 500 {
			lg.Error("http request")
		} else {
			lg.Debug("http request")
		}
	}
}
