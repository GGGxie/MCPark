package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func limitRouter(max int64) *gin.Engine {
	r := gin.New()
	r.Use(MaxBodyBytes(max))
	r.POST("/echo", func(c *gin.Context) {
		b, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusRequestEntityTooLarge)
			return
		}
		c.String(http.StatusOK, "%d", len(b))
	})
	return r
}

func TestMaxBodyBytesAllowsSmallBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(strings.Repeat("a", 100)))
	w := httptest.NewRecorder()
	limitRouter(100).ServeHTTP(w, req)
	if w.Code != http.StatusOK || w.Body.String() != "100" {
		t.Fatalf("got %d %q", w.Code, w.Body.String())
	}
}

func TestMaxBodyBytesRejectsDeclaredOversize(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(strings.Repeat("a", 101)))
	w := httptest.NewRecorder()
	limitRouter(100).ServeHTTP(w, req)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", w.Code)
	}
}

func TestMaxBodyBytesCutsOffUndeclaredOversize(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(strings.Repeat("a", 5000)))
	req.ContentLength = -1 // e.g. chunked transfer encoding
	w := httptest.NewRecorder()
	limitRouter(100).ServeHTTP(w, req)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", w.Code)
	}
}
