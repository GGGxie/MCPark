package handler

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

// fakeGoogle stands in for Google's token and userinfo endpoints.
type fakeGoogle struct {
	srv          *httptest.Server
	tokenCalls   atomic.Int32
	userinfoJSON string
}

func newFakeGoogle(t *testing.T) *fakeGoogle {
	t.Helper()
	f := &fakeGoogle{userinfoJSON: `{"id":"google-123","email":"someone@example.com","name":"Someone","picture":"p"}`}
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		f.tokenCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token":"at","token_type":"Bearer","expires_in":3600}`))
	})
	mux.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(f.userinfoJSON))
	})
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

func newAuthRouter(t *testing.T, frontendURL string) (*gin.Engine, *AuthHandler, *fakeGoogle) {
	t.Helper()
	cfg := testConfig()
	cfg.FrontendURL = frontendURL
	h := NewAuthHandler(newTestDB(t), cfg)
	g := newFakeGoogle(t)
	h.oauthConfig.Endpoint = oauth2.Endpoint{AuthURL: "https://accounts.example.com/auth", TokenURL: g.srv.URL + "/token"}
	h.userInfoURL = g.srv.URL + "/userinfo"

	r := gin.New()
	r.GET("/api/auth/google", h.GoogleLogin)
	r.GET("/api/auth/google/callback", h.GoogleCallback)
	return r, h, g
}

func findCookie(w *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, c := range w.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func callbackReq(state string, cookieValue *string) *http.Request {
	target := "/api/auth/google/callback?code=abc"
	if state != "" {
		target += "&state=" + url.QueryEscape(state)
	}
	req := httptest.NewRequest(http.MethodGet, target, nil)
	if cookieValue != nil {
		req.AddCookie(&http.Cookie{Name: oauthStateCookie, Value: *cookieValue})
	}
	return req
}

func TestGoogleLoginSetsRandomStateCookie(t *testing.T) {
	r, _, _ := newAuthRouter(t, "https://app.example.com")

	w1 := serve(r, httptest.NewRequest(http.MethodGet, "/api/auth/google", nil))
	if w1.Code != http.StatusOK {
		t.Fatalf("status = %d", w1.Code)
	}
	c := findCookie(w1, oauthStateCookie)
	if c == nil {
		t.Fatal("state cookie not set")
	}
	raw, err := base64.RawURLEncoding.DecodeString(c.Value)
	if err != nil || len(raw) != 32 {
		t.Fatalf("state must be 32 random bytes in base64url, got len=%d err=%v", len(raw), err)
	}
	if !c.HttpOnly || c.SameSite != http.SameSiteLaxMode || c.Path != "/api/auth" || c.MaxAge != 600 || !c.Secure {
		t.Fatalf("unexpected cookie attributes: %+v", c)
	}
	if !strings.Contains(w1.Body.String(), "state="+c.Value) {
		t.Fatalf("authorization URL does not carry the state: %s", w1.Body.String())
	}

	w2 := serve(r, httptest.NewRequest(http.MethodGet, "/api/auth/google", nil))
	if c2 := findCookie(w2, oauthStateCookie); c2 == nil || c2.Value == c.Value {
		t.Fatal("state must differ between requests")
	}
}

func TestGoogleLoginCookieNotSecureOnHTTPFrontend(t *testing.T) {
	r, _, _ := newAuthRouter(t, "http://localhost:5173")
	c := findCookie(serve(r, httptest.NewRequest(http.MethodGet, "/api/auth/google", nil)), oauthStateCookie)
	if c == nil || c.Secure {
		t.Fatalf("cookie must exist and must not be Secure for http frontend: %+v", c)
	}
}

func TestCallbackRejectsBadState(t *testing.T) {
	good := "abcDEF123_-abcDEF123_-abcDEF123_-abcDEF123_-"
	other := "zzzzzz123_-abcDEF123_-abcDEF123_-abcDEF123_-"
	empty := ""
	cases := map[string]*http.Request{
		"no cookie, no state":    callbackReq("", nil),
		"cookie but no state":    callbackReq("", &good),
		"state but no cookie":    callbackReq(good, nil),
		"mismatch":               callbackReq(other, &good),
		"empty cookie and state": callbackReq("", &empty),
	}
	for name, req := range cases {
		t.Run(name, func(t *testing.T) {
			r, _, g := newAuthRouter(t, "https://app.example.com")
			w := serve(r, req)
			if w.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", w.Code)
			}
			c := findCookie(w, oauthStateCookie)
			if c == nil || c.MaxAge >= 0 {
				t.Fatalf("state cookie must be cleared, got %+v", c)
			}
			if n := g.tokenCalls.Load(); n != 0 {
				t.Fatalf("token endpoint must not be reached on bad state (calls=%d)", n)
			}
		})
	}
}

func TestCallbackSuccessPutsTokenInFragment(t *testing.T) {
	r, _, g := newAuthRouter(t, "https://app.example.com")
	state := "validstate_validstate_validstate_valid"
	w := serve(r, callbackReq(state, &state))

	if w.Code != http.StatusTemporaryRedirect && w.Code != http.StatusFound {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if g.tokenCalls.Load() != 1 {
		t.Fatal("expected one token exchange")
	}
	loc := w.Header().Get("Location")
	const prefix = "https://app.example.com/auth/callback#token="
	if !strings.HasPrefix(loc, prefix) {
		t.Fatalf("Location = %q, want prefix %q", loc, prefix)
	}
	if strings.Contains(loc, "?") {
		t.Fatalf("token must not travel in a query string: %q", loc)
	}
	parsed, err := jwt.Parse(strings.TrimPrefix(loc, prefix), func(*jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())
	if err != nil || !parsed.Valid {
		t.Fatalf("issued token invalid: %v", err)
	}
	if c := findCookie(w, oauthStateCookie); c == nil || c.MaxAge >= 0 {
		t.Fatalf("state cookie should be cleared after use, got %+v", c)
	}
	if cc := w.Header().Get("Cache-Control"); !strings.Contains(cc, "no-store") {
		t.Fatalf("redirect should not be cacheable, Cache-Control=%q", cc)
	}
}
