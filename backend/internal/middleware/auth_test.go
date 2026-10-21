package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "0123456789abcdef0123456789abcdef-test"

func init() { gin.SetMode(gin.TestMode) }

func newRouter() *gin.Engine {
	r := gin.New()
	r.GET("/api/me", AuthRequired(testSecret), func(c *gin.Context) {
		uid, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"user_id": uid})
	})
	r.GET("/api/admin/ping", AdminRequired(testSecret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func signHS256(t *testing.T, secret string, claims jwt.MapClaims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func validClaims(isAdmin bool) jwt.MapClaims {
	return jwt.MapClaims{
		"user_id":  1,
		"email":    "u@example.com",
		"name":     "u",
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
}

func do(r http.Handler, path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthRequiredAcceptsValidToken(t *testing.T) {
	w := do(newRouter(), "/api/me", signHS256(t, testSecret, validClaims(false)))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestRejectedTokens(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	rs256, err := jwt.NewWithClaims(jwt.SigningMethodRS256, validClaims(true)).SignedString(rsaKey)
	if err != nil {
		t.Fatal(err)
	}
	none, err := jwt.NewWithClaims(jwt.SigningMethodNone, validClaims(true)).SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatal(err)
	}
	hs384, err := jwt.NewWithClaims(jwt.SigningMethodHS384, validClaims(true)).SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	valid := signHS256(t, testSecret, validClaims(true))
	parts := strings.Split(valid, ".")
	tampered := parts[0] + "." + parts[1] + "." + parts[2][:len(parts[2])-2] + "xx"

	expired := validClaims(true)
	expired["exp"] = time.Now().Add(-time.Hour).Unix()

	noExp := validClaims(true)
	delete(noExp, "exp")

	noUserID := validClaims(false)
	delete(noUserID, "user_id")

	cases := map[string]string{
		"missing":          "",
		"garbage":          "not-a-jwt",
		"tampered":         tampered,
		"wrong secret":     signHS256(t, "another-secret-another-secret-another", validClaims(true)),
		"expired":          signHS256(t, testSecret, expired),
		"no exp":           signHS256(t, testSecret, noExp),
		"no user_id":       signHS256(t, testSecret, noUserID),
		"alg none":         none,
		"alg RS256":        rs256,
		"unexpected HS384": hs384,
	}
	r := newRouter()
	for name, tok := range cases {
		for _, path := range []string{"/api/me", "/api/admin/ping"} {
			t.Run(name+" "+path, func(t *testing.T) {
				if w := do(r, path, tok); w.Code != http.StatusUnauthorized {
					t.Fatalf("status = %d, want 401", w.Code)
				}
			})
		}
	}
}

func TestAdminRouteRequiresAdminClaim(t *testing.T) {
	r := newRouter()
	if w := do(r, "/api/admin/ping", signHS256(t, testSecret, validClaims(false))); w.Code != http.StatusForbidden {
		t.Fatalf("non-admin status = %d, want 403", w.Code)
	}
	if w := do(r, "/api/admin/ping", signHS256(t, testSecret, validClaims(true))); w.Code != http.StatusOK {
		t.Fatalf("admin status = %d, want 200", w.Code)
	}
}

func TestAdminClaimWrongTypeIsNotAdmin(t *testing.T) {
	c := validClaims(false)
	c["is_admin"] = "true" // string, not bool
	w := do(newRouter(), "/api/admin/ping", signHS256(t, testSecret, c))
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestMalformedClaimTypesDoNotPanic(t *testing.T) {
	c := jwt.MapClaims{
		"user_id":  "not-a-number",
		"email":    123,
		"name":     []string{"x"},
		"is_admin": 1,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}
	r := newRouter()
	tok := signHS256(t, testSecret, c)
	if w := do(r, "/api/me", tok); w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
	if w := do(r, "/api/admin/ping", tok); w.Code != http.StatusUnauthorized {
		t.Fatalf("admin status = %d, want 401", w.Code)
	}
}

func TestAuthOptionalIgnoresBadToken(t *testing.T) {
	r := gin.New()
	r.GET("/x", AuthOptional(testSecret), func(c *gin.Context) {
		_, ok := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{"auth": ok})
	})
	if w := do(r, "/x", "garbage"); w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"auth":false`) {
		t.Fatalf("unexpected: %d %s", w.Code, w.Body.String())
	}
	if w := do(r, "/x", signHS256(t, testSecret, validClaims(false))); !strings.Contains(w.Body.String(), `"auth":true`) {
		t.Fatalf("expected authenticated: %s", w.Body.String())
	}
}
