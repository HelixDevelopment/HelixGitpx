package gin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	hgin "github.com/helixgitpx/platform/gin"
)

func TestNewRouter_BaseHeaders(t *testing.T) {
	r := hgin.NewRouter(hgin.Options{Service: "hello", Version: "test"})
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("code = %d", w.Code)
	}
	if w.Header().Get("X-HelixGitpx-Service") != "hello" {
		t.Errorf("missing service header")
	}
	if w.Header().Get("X-HelixGitpx-Version") != "test" {
		t.Errorf("missing version header")
	}
}
