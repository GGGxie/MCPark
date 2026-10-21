package handler

import (
	"fmt"
	"net/http"
	"testing"
)

func TestDeleteServiceRemovesDependentRows(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	uid := insertUser(t, database, "u@example.com", false)
	admin := signToken(t, insertUser(t, database, "a@example.com", true), true)
	sid := insertService(t, database, "approved")

	database.Exec("INSERT INTO likes (user_id, service_id) VALUES (?, ?)", uid, sid)
	database.Exec("INSERT INTO comments (service_id, user_id, content) VALUES (?, ?, 'c')", sid, uid)

	w := call(t, r, http.MethodDelete, fmt.Sprintf("/api/admin/services/%d", sid), admin, "", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	for _, q := range []string{
		"SELECT COUNT(*) FROM services WHERE id = ?",
		"SELECT COUNT(*) FROM likes WHERE service_id = ?",
		"SELECT COUNT(*) FROM comments WHERE service_id = ?",
	} {
		if n := scanInt(t, database, q, sid); n != 0 {
			t.Fatalf("%s -> %d, want 0", q, n)
		}
	}
}

func TestDeleteServiceIsAtomic(t *testing.T) {
	database := newTestDB(t)
	r := newTestRouter(database)
	uid := insertUser(t, database, "u@example.com", false)
	admin := signToken(t, insertUser(t, database, "a@example.com", true), true)
	sid := insertService(t, database, "approved")

	database.Exec("INSERT INTO likes (user_id, service_id) VALUES (?, ?)", uid, sid)
	database.Exec("INSERT INTO comments (service_id, user_id, content) VALUES (?, ?, 'c')", sid, uid)

	// Make the final DELETE fail: the earlier deletes must be rolled back.
	if _, err := database.Exec(`CREATE TRIGGER fail_service_delete BEFORE DELETE ON services
		BEGIN SELECT RAISE(ABORT, 'forced failure'); END`); err != nil {
		t.Fatal(err)
	}
	w := call(t, r, http.MethodDelete, fmt.Sprintf("/api/admin/services/%d", sid), admin, "", nil)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if n := scanInt(t, database, "SELECT COUNT(*) FROM likes WHERE service_id = ?", sid); n != 1 {
		t.Fatalf("likes were not rolled back (count=%d)", n)
	}
	if n := scanInt(t, database, "SELECT COUNT(*) FROM comments WHERE service_id = ?", sid); n != 1 {
		t.Fatalf("comments were not rolled back (count=%d)", n)
	}
}
