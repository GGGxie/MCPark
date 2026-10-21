package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"MCPark/internal/config"
	"MCPark/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testJWTSecret = "unit-test-secret-unit-test-secret-0123456789"

func init() { gin.SetMode(gin.TestMode) }

func testConfig() *config.Config {
	return &config.Config{
		JWTSecret:         testJWTSecret,
		FrontendURL:       "https://app.example.com",
		AdminEmail:        "admin@example.com",
		GoogleClientID:    "test-client-id",
		GoogleRedirectURL: "https://app.example.com/api/auth/google/callback",
	}
}

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	database, err := db.Init(db.MemoryDSN(t.Name()))
	if err != nil {
		t.Fatalf("init test db: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	return database
}

func signToken(t *testing.T, userID int64, isAdmin bool) string {
	t.Helper()
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  userID,
		"email":    "user@example.com",
		"name":     "Test User",
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte(testJWTSecret))
	if err != nil {
		t.Fatal(err)
	}
	return tok
}

func insertUser(t *testing.T, database *sql.DB, email string, admin bool) int64 {
	t.Helper()
	a := 0
	if admin {
		a = 1
	}
	res, err := database.Exec("INSERT INTO users (google_id, email, name, is_admin) VALUES (?, ?, ?, ?)", "g-"+email, email, "Test User", a)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return id
}

func insertService(t *testing.T, database *sql.DB, status string) int64 {
	t.Helper()
	res, err := database.Exec(`INSERT INTO services (name, title, description, category, author_id, author_name, status)
		VALUES ('n', 'Test Service', 'desc', 'Testing', 1, 'MCPark', ?)`, status)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return id
}

func jsonBody(t *testing.T, v any) *bytes.Reader {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(b)
}

func serve(r http.Handler, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
