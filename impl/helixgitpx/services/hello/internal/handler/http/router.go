package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/helixgitpx/helixgitpx/services/hello/internal/domain"
)

// Register adds /v1/hello to r.
func Register(r *gin.Engine, g *domain.Greeter) {
	r.GET("/v1/hello", func(c *gin.Context) {
		name := c.Query("name")
		resp, err := g.Greet(c.Request.Context(), name)
		if err != nil {
			c.JSON(nethttp.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(nethttp.StatusOK, gin.H{"greeting": resp.Greeting, "count": resp.Count})
	})
}
