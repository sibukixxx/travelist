package sqlite_test

import (
	"testing"

	"github.com/sibukixxx/travelist/api/internal/infra/repo/sqlite"
)

func TestOpenInMemory(t *testing.T) {
	db, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("Open(:memory:) error = %v", err)
	}
	defer db.Close()

	// Verify the users table was created by the migration
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableName)
	if err != nil {
		t.Fatalf("expected users table to exist: %v", err)
	}
	if tableName != "users" {
		t.Errorf("table name = %q, want %q", tableName, "users")
	}
}

func TestOpenRunsMigrationsIdempotently(t *testing.T) {
	db, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("first Open error = %v", err)
	}
	defer db.Close()

	// Running Migrate again should not fail
	if err := sqlite.Migrate(db); err != nil {
		t.Fatalf("second Migrate error = %v", err)
	}
}
