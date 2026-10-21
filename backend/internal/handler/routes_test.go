package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestHealthEndpointIsPublic(t *testing.T) {
	r := newTestRouter(newTestDB(t))
	w := call(t, r, http.MethodGet, "/api/health", "", "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil || body["status"] != "ok" || len(body) != 1 {
		t.Fatalf("body = %s (err=%v)", w.Body.String(), err)
	}
}

func TestAdminRoutesAuthorization(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	user := signToken(t, insertUser(t, database, "u@example.com", false), false)
	admin := signToken(t, insertUser(t, database, "a@example.com", true), true)

	routes := []struct{ method, path string }{
		{http.MethodGet, "/api/admin/services"},
		{http.MethodPost, "/api/admin/services/1/approve"},
		{http.MethodPost, "/api/admin/services/1/reject"},
		{http.MethodDelete, "/api/admin/services/1"},
	}
	for _, rt := range routes {
		t.Run(rt.method+" "+rt.path, func(t *testing.T) {
			if w := call(t, r, rt.method, rt.path, "", "", nil); w.Code != http.StatusUnauthorized {
				t.Errorf("anonymous: %d, want 401", w.Code)
			}
			if w := call(t, r, rt.method, rt.path, user, "", nil); w.Code != http.StatusForbidden {
				t.Errorf("non-admin: %d, want 403", w.Code)
			}
			if w := call(t, r, rt.method, rt.path, admin, "", nil); w.Code != http.StatusOK {
				t.Errorf("admin: %d, want 200", w.Code)
			}
		})
	}
}

func TestOversizedBodyIsRejected(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	tok := signToken(t, insertUser(t, database, "u@example.com", false), false)

	huge := `{"title":"t","description":"d","category":"c","content":"` + strings.Repeat("a", 300<<10) + `"}`

	// Declared length.
	if w := call(t, r, http.MethodPost, "/api/services", tok, "", bytes.NewReader([]byte(huge))); w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("declared oversize: status = %d, want 413", w.Code)
	}

	// Undeclared length (chunked): the handler must still refuse with 413.
	req, _ := http.NewRequest(http.MethodPost, "/api/services", strings.NewReader(huge))
	req.ContentLength = -1
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	if w := serve(r, req); w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("chunked oversize: status = %d, want 413 (body=%s)", w.Code, w.Body.String())
	}

	if n := scanInt(t, database, "SELECT COUNT(*) FROM services WHERE title = 't'"); n != 0 {
		t.Fatalf("oversized request must not create a service")
	}
}
