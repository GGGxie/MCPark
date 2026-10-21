package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"MCPark/internal/config"
	"MCPark/internal/db"

	"github.com/gin-gonic/gin"
)

func testCfg() *config.Config {
	return &config.Config{
		Port:        "8086",
		ListenAddr:  "127.0.0.1",
		JWTSecret:   strings.Repeat("k", 40),
		FrontendURL: "http://localhost:5173",
	}
}

func TestResolveGinMode(t *testing.T) {
	if got := resolveGinMode(func(string) string { return "" }); got != gin.ReleaseMode {
		t.Errorf("default mode = %q, want release", got)
	}
	if got := resolveGinMode(func(k string) string {
		if k == "GIN_MODE" {
			return "debug"
		}
		return ""
	}); got != gin.DebugMode {
		t.Errorf("explicit GIN_MODE ignored, got %q", got)
	}
}

func TestListenAddress(t *testing.T) {
	cfg := testCfg()
	if got := listenAddress(cfg); got != "127.0.0.1:8086" {
		t.Errorf("default = %q", got)
	}
	cfg.ListenAddr = "0.0.0.0"
	if got := listenAddress(cfg); got != "0.0.0.0:8086" {
		t.Errorf("explicit = %q", got)
	}
	cfg.ListenAddr = "::1"
	if got := listenAddress(cfg); got != "[::1]:8086" {
		t.Errorf("ipv6 = %q", got)
	}
}

func newTestEngine(t *testing.T) *gin.Engine {
	t.Helper()
	database, err := db.Init(db.MemoryDSN(t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })
	r, err := newEngine(testCfg(), database)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestEngineOnlyTrustsLoopbackProxy(t *testing.T) {
	r := newTestEngine(t)
	var seen string
	r.GET("/ip", func(c *gin.Context) { seen = c.ClientIP() })

	for _, tc := range []struct{ remote, want string }{
		{"203.0.113.5:1234", "203.0.113.5"}, // untrusted peer: forwarded header ignored
		{"127.0.0.1:1234", "198.51.100.7"},  // local reverse proxy: forwarded header honoured
	} {
		req := httptest.NewRequest(http.MethodGet, "/ip", nil)
		req.RemoteAddr = tc.remote
		req.Header.Set("X-Forwarded-For", "198.51.100.7")
		r.ServeHTTP(httptest.NewRecorder(), req)
		if seen != tc.want {
			t.Errorf("remote %s: ClientIP = %q, want %q", tc.remote, seen, tc.want)
		}
	}
}

func TestEngineLimitsBodyAndServesHealth(t *testing.T) {
	r := newTestEngine(t)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/health", nil))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Fatalf("health: %d %s", w.Code, w.Body.String())
	}

	big := strings.NewReader(`{"credential":"` + strings.Repeat("a", 300<<10) + `"}`)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/api/auth/google/token", big))
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("oversize body: status = %d, want 413", w.Code)
	}
}

func TestUnknownAPIPathReturnsJSON404(t *testing.T) {
	r := newTestEngine(t)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/nope", nil))
	if w.Code != http.StatusNotFound || !strings.Contains(w.Body.String(), "Not found") {
		t.Fatalf("got %d %s", w.Code, w.Body.String())
	}
}
