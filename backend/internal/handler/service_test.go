package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func validServiceReq() map[string]any {
	return map[string]any{
		"name":        "my-server",
		"title":       "My Server",
		"description": "does things",
		"category":    "Testing",
		"content":     "# Hello",
	}
}

func TestCreateServiceAcceptsValidInputAsPending(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	uid := insertUser(t, database, "u@example.com", false)

	w := call(t, r, http.MethodPost, "/api/services", signToken(t, uid, false), "", validServiceReq())
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if n := scanInt(t, database, "SELECT COUNT(*) FROM services WHERE title = 'My Server' AND status = 'pending'"); n != 1 {
		t.Fatalf("expected one pending service, got %d", n)
	}
}

func TestCreateServiceRequiresAuth(t *testing.T) {
	r := newTestRouter(newTestDB(t))
	if w := call(t, r, http.MethodPost, "/api/services", "", "", validServiceReq()); w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

func TestCreateServiceFieldLengthLimits(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	tok := signToken(t, insertUser(t, database, "u@example.com", false), false)

	limits := map[string]int{
		"name":        200,
		"title":       200,
		"title_zh":    200,
		"description": 2000,
		"desc_zh":     2000,
		"content":     50000,
		"content_zh":  50000,
		"config_json": 50000,
	}
	for field, max := range limits {
		t.Run(field+" at limit", func(t *testing.T) {
			req := validServiceReq()
			req[field] = strings.Repeat("a", max)
			if w := call(t, r, http.MethodPost, "/api/services", tok, "", req); w.Code != http.StatusCreated {
				t.Fatalf("at limit: status = %d body=%.200s", w.Code, w.Body.String())
			}
		})
		t.Run(field+" over limit", func(t *testing.T) {
			req := validServiceReq()
			req[field] = strings.Repeat("a", max+1)
			if w := call(t, r, http.MethodPost, "/api/services", tok, "", req); w.Code != http.StatusBadRequest {
				t.Fatalf("over limit: status = %d, want 400", w.Code)
			}
		})
	}
}

func TestCreateCommentLengthLimit(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	tok := signToken(t, insertUser(t, database, "u@example.com", false), false)
	sid := insertService(t, database, "approved")
	path := fmt.Sprintf("/api/services/%d/comments", sid)

	if w := call(t, r, http.MethodPost, path, tok, "", map[string]string{"content": strings.Repeat("a", 2000)}); w.Code != http.StatusCreated {
		t.Fatalf("at limit: status = %d body=%s", w.Code, w.Body.String())
	}
	if w := call(t, r, http.MethodPost, path, tok, "", map[string]string{"content": strings.Repeat("a", 2001)}); w.Code != http.StatusBadRequest {
		t.Fatalf("over limit: status = %d, want 400", w.Code)
	}
	if w := call(t, r, http.MethodPost, path, tok, "", map[string]string{"content": ""}); w.Code != http.StatusBadRequest {
		t.Fatalf("empty: status = %d, want 400", w.Code)
	}
}

func TestCreateCommentRequiresApprovedService(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	tok := signToken(t, insertUser(t, database, "u@example.com", false), false)
	body := map[string]string{"content": "hi"}

	cases := map[string]string{
		"missing":  "/api/services/99999/comments",
		"pending":  fmt.Sprintf("/api/services/%d/comments", insertService(t, database, "pending")),
		"rejected": fmt.Sprintf("/api/services/%d/comments", insertService(t, database, "rejected")),
		"non-id":   "/api/services/abc/comments",
	}
	for name, path := range cases {
		t.Run(name, func(t *testing.T) {
			if w := call(t, r, http.MethodPost, path, tok, "", body); w.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", w.Code)
			}
		})
	}
	if n := scanInt(t, database, "SELECT COUNT(*) FROM comments"); n != 0 {
		t.Fatalf("no comment should be stored, got %d", n)
	}
}

func TestLikeToggleOnApprovedService(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	uid := insertUser(t, database, "u@example.com", false)
	tok := signToken(t, uid, false)
	sid := insertService(t, database, "approved")
	path := fmt.Sprintf("/api/services/%d/like", sid)

	w := call(t, r, http.MethodPost, path, tok, "", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"liked":true`) {
		t.Fatalf("like: %d %s", w.Code, w.Body.String())
	}
	if n := scanInt(t, database, "SELECT likes FROM services WHERE id = ?", sid); n != 1 {
		t.Fatalf("likes = %d, want 1", n)
	}

	w = call(t, r, http.MethodPost, path, tok, "", nil)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"liked":false`) {
		t.Fatalf("unlike: %d %s", w.Code, w.Body.String())
	}
	if n := scanInt(t, database, "SELECT likes FROM services WHERE id = ?", sid); n != 0 {
		t.Fatalf("likes = %d, want 0", n)
	}
	if n := scanInt(t, database, "SELECT COUNT(*) FROM likes"); n != 0 {
		t.Fatalf("likes rows = %d, want 0", n)
	}
}

func TestLikeRequiresApprovedService(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	tok := signToken(t, insertUser(t, database, "u@example.com", false), false)

	for name, path := range map[string]string{
		"missing": "/api/services/99999/like",
		"pending": fmt.Sprintf("/api/services/%d/like", insertService(t, database, "pending")),
		"non-id":  "/api/services/abc/like",
	} {
		t.Run(name, func(t *testing.T) {
			if w := call(t, r, http.MethodPost, path, tok, "", nil); w.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", w.Code)
			}
		})
	}
	if n := scanInt(t, database, "SELECT COUNT(*) FROM likes"); n != 0 {
		t.Fatalf("no like row expected, got %d", n)
	}
}

func TestLikeIsAtomicWhenCounterUpdateFails(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	tok := signToken(t, insertUser(t, database, "u@example.com", false), false)
	sid := insertService(t, database, "approved")

	// Force the counter UPDATE to fail; the like row inserted earlier in the
	// same operation must be rolled back and the client must see an error.
	if _, err := database.Exec(`CREATE TRIGGER fail_like_update BEFORE UPDATE OF likes ON services
		BEGIN SELECT RAISE(ABORT, 'forced failure'); END`); err != nil {
		t.Fatal(err)
	}
	w := call(t, r, http.MethodPost, fmt.Sprintf("/api/services/%d/like", sid), tok, "", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if n := scanInt(t, database, "SELECT COUNT(*) FROM likes"); n != 0 {
		t.Fatalf("like row was not rolled back (count=%d)", n)
	}
}

func TestListServicesIsSafeAgainstInjectionStyleInput(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	total := scanInt(t, database, "SELECT COUNT(*) FROM services WHERE status = 'approved'")
	if total == 0 {
		t.Fatal("seed data expected")
	}

	payloads := []string{
		"' OR '1'='1",
		"' OR 1=1 --",
		"'; DROP TABLE services; --",
		`" OR ""="`,
		"%') UNION SELECT 1,2,3 --",
	}
	for _, p := range payloads {
		for _, param := range []string{"search", "category"} {
			t.Run(param+"="+p, func(t *testing.T) {
				w := call(t, r, http.MethodGet, "/api/services?"+param+"="+urlEscape(p), "", "", nil)
				if w.Code != http.StatusOK {
					t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
				}
				var resp struct {
					Total int               `json:"total"`
					Data  []json.RawMessage `json:"data"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatal(err)
				}
				if resp.Total != 0 || len(resp.Data) != 0 {
					t.Fatalf("payload matched rows as data: total=%d len=%d", resp.Total, len(resp.Data))
				}
			})
		}
	}
	if n := scanInt(t, database, "SELECT COUNT(*) FROM services WHERE status = 'approved'"); n != total {
		t.Fatalf("services table changed: %d -> %d", total, n)
	}
}

func urlEscape(s string) string {
	return strings.NewReplacer(" ", "%20", "'", "%27", "\"", "%22", "%", "%25", ";", "%3B", "#", "%23", "&", "%26", "+", "%2B", "=", "%3D", ",", "%2C", "(", "%28", ")", "%29").Replace(s)
}

func TestListServicesPagingIsClamped(t *testing.T) {
	r := newTestRouter(newTestDB(t))
	cases := []struct {
		query    string
		wantPage int
		wantSize int
	}{
		{"", 1, 12},
		{"?page_size=1000", 1, 50},
		{"?page_size=50", 1, 50},
		{"?page_size=0", 1, 12},
		{"?page_size=-3", 1, 12},
		{"?page_size=abc", 1, 12},
		{"?page=-5", 1, 12},
		{"?page=0", 1, 12},
		{"?page=9223372036854775807&page_size=50", maxPage, 50},
	}
	for _, tc := range cases {
		t.Run(tc.query, func(t *testing.T) {
			w := call(t, r, http.MethodGet, "/api/services"+tc.query, "", "", nil)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
			}
			var resp struct {
				Page     int `json:"page"`
				PageSize int `json:"page_size"`
			}
			json.Unmarshal(w.Body.Bytes(), &resp)
			if resp.Page != tc.wantPage || resp.PageSize != tc.wantSize {
				t.Fatalf("page=%d size=%d, want page=%d size=%d", resp.Page, resp.PageSize, tc.wantPage, tc.wantSize)
			}
		})
	}
}

func TestAdminListPagingIsClamped(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	admin := signToken(t, insertUser(t, database, "a@example.com", true), true)
	w := call(t, r, http.MethodGet, "/api/admin/services?page_size=100000", admin, "", nil)
	var resp struct {
		PageSize int `json:"page_size"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if w.Code != http.StatusOK || resp.PageSize != 100 {
		t.Fatalf("status=%d page_size=%d, want 200/100", w.Code, resp.PageSize)
	}
}

func TestInstallCountIsDeduplicatedPerClient(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	sid := insertService(t, database, "approved")
	path := fmt.Sprintf("/api/services/%d/install", sid)

	w := call(t, r, http.MethodPost, path, "", "10.0.0.1", nil)
	if !strings.Contains(w.Body.String(), `"counted":true`) {
		t.Fatalf("first install should count: %s", w.Body.String())
	}
	w = call(t, r, http.MethodPost, path, "", "10.0.0.1", nil)
	if !strings.Contains(w.Body.String(), `"counted":false`) {
		t.Fatalf("repeat install should not count: %s", w.Body.String())
	}
	if n := scanInt(t, database, "SELECT install_count FROM services WHERE id = ?", sid); n != 1 {
		t.Fatalf("install_count = %d, want 1", n)
	}
	call(t, r, http.MethodPost, path, "", "10.0.0.2", nil)
	if n := scanInt(t, database, "SELECT install_count FROM services WHERE id = ?", sid); n != 2 {
		t.Fatalf("a different client should count: install_count = %d, want 2", n)
	}
}

func TestInstallRequiresApprovedService(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	pending := insertService(t, database, "pending")
	for name, path := range map[string]string{
		"missing": "/api/services/99999/install",
		"pending": fmt.Sprintf("/api/services/%d/install", pending),
	} {
		t.Run(name, func(t *testing.T) {
			if w := call(t, r, http.MethodPost, path, "", "10.0.0.9", nil); w.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404", w.Code)
			}
		})
	}
	if n := scanInt(t, database, "SELECT install_count FROM services WHERE id = ?", pending); n != 0 {
		t.Fatalf("install_count changed for unapproved service: %d", n)
	}
}

func TestViewCountIsDeduplicatedPerClient(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	sid := insertService(t, database, "approved")
	path := fmt.Sprintf("/api/services/%d", sid)

	for i := 0; i < 3; i++ {
		if w := call(t, r, http.MethodGet, path, "", "10.0.1.1", nil); w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	}
	if n := scanInt(t, database, "SELECT views FROM services WHERE id = ?", sid); n != 1 {
		t.Fatalf("views = %d, want 1", n)
	}
	call(t, r, http.MethodGet, path, "", "10.0.1.2", nil)
	if n := scanInt(t, database, "SELECT views FROM services WHERE id = ?", sid); n != 2 {
		t.Fatalf("views = %d, want 2", n)
	}
}

func TestGetServiceHidesUnapproved(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	pending := insertService(t, database, "pending")
	if w := call(t, r, http.MethodGet, fmt.Sprintf("/api/services/%d", pending), "", "10.0.2.1", nil); w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if n := scanInt(t, database, "SELECT views FROM services WHERE id = ?", pending); n != 0 {
		t.Fatalf("views changed for unapproved service: %d", n)
	}
}
