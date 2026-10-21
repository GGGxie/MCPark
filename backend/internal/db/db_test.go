package db

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestInitMemorySchemaAndSeedWithForeignKeys(t *testing.T) {
	database, err := Init(MemoryDSN(t.Name()))
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	t.Cleanup(func() { database.Close() })

	var fk int
	if err := database.QueryRow("PRAGMA foreign_keys").Scan(&fk); err != nil || fk != 1 {
		t.Fatalf("foreign_keys = %d, err=%v; want 1", fk, err)
	}

	var services, articles int
	database.QueryRow("SELECT COUNT(*) FROM services").Scan(&services)
	database.QueryRow("SELECT COUNT(*) FROM articles").Scan(&articles)
	if services == 0 || articles == 0 {
		t.Fatalf("seed incomplete: services=%d articles=%d", services, articles)
	}

	// Seeded services must all reference an existing author.
	var orphans int
	database.QueryRow("SELECT COUNT(*) FROM services s LEFT JOIN users u ON u.id = s.author_id WHERE u.id IS NULL").Scan(&orphans)
	if orphans != 0 {
		t.Fatalf("%d seeded services without author", orphans)
	}
}

func TestForeignKeysAreEnforced(t *testing.T) {
	database, err := Init(MemoryDSN(t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { database.Close() })

	_, err = database.Exec("INSERT INTO comments (service_id, user_id, content) VALUES (99999, 1, 'x')")
	if err == nil {
		t.Fatal("expected foreign key violation for unknown service_id")
	}
}

func TestInitIsIdempotent(t *testing.T) {
	dsn := MemoryDSN(t.Name())
	first, err := Init(dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { first.Close() })
	var before int
	first.QueryRow("SELECT COUNT(*) FROM services").Scan(&before)

	// A second initialisation against the same database must not duplicate seed data.
	second, err := Init(dsn)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { second.Close() })
	var after int
	second.QueryRow("SELECT COUNT(*) FROM services").Scan(&after)
	if before != after {
		t.Fatalf("services count changed %d -> %d", before, after)
	}
}

func TestFileDSNEnablesForeignKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "t.db")
	dsn := FileDSN(path)
	if !strings.Contains(dsn, "_foreign_keys=on") {
		t.Fatalf("file DSN lacks foreign key option: %s", dsn)
	}
	database := InitDB(path)
	defer database.Close()
	var fk int
	if err := database.QueryRow("PRAGMA foreign_keys").Scan(&fk); err != nil || fk != 1 {
		t.Fatalf("foreign_keys = %d, err=%v", fk, err)
	}
}
