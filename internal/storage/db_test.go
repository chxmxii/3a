package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenCreatesDBAndSchema(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "sub", "nested", "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}
	defer store.Close()

	// Verify the file was created (including parent directories).
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("expected database file to be created")
	}

	// Verify all 6 tables exist by querying sqlite_master.
	expectedTables := []string{
		"assessments",
		"resources",
		"relationships",
		"findings",
		"cost_estimates",
		"sizing",
	}

	for _, table := range expectedTables {
		var name string
		err := store.DB.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='table' AND name=?", table,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %q not found: %v", table, err)
		}
	}
}

func TestOpenCreatesIndexes(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}
	defer store.Close()

	expectedIndexes := []string{
		"idx_resources_assessment",
		"idx_resources_type",
		"idx_resources_region",
		"idx_relationships_assessment",
		"idx_relationships_source",
		"idx_findings_assessment",
		"idx_findings_severity",
		"idx_findings_category",
		"idx_costs_assessment",
		"idx_sizing_assessment",
		"idx_sizing_category",
	}

	for _, idx := range expectedIndexes {
		var name string
		err := store.DB.QueryRow(
			"SELECT name FROM sqlite_master WHERE type='index' AND name=?", idx,
		).Scan(&name)
		if err != nil {
			t.Errorf("index %q not found: %v", idx, err)
		}
	}
}

func TestOpenIdempotent(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	// Open twice to verify migrations are idempotent (IF NOT EXISTS).
	store1, err := Open(dbPath)
	if err != nil {
		t.Fatalf("first Open() error: %v", err)
	}
	store1.Close()

	store2, err := Open(dbPath)
	if err != nil {
		t.Fatalf("second Open() error: %v", err)
	}
	defer store2.Close()

	// Verify tables still intact after second Open.
	var count int
	err = store2.DB.QueryRow(
		"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'",
	).Scan(&count)
	if err != nil {
		t.Fatalf("querying tables: %v", err)
	}
	if count != 6 {
		t.Errorf("expected 6 tables, got %d", count)
	}
}

func TestStoreClose(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error: %v", err)
	}

	// After closing, queries should fail.
	err = store.DB.Ping()
	if err == nil {
		t.Error("expected error after Close(), got nil")
	}
}

func TestOpenForeignKeysEnabled(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	store, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}
	defer store.Close()

	var fkEnabled int
	err = store.DB.QueryRow("PRAGMA foreign_keys").Scan(&fkEnabled)
	if err != nil {
		t.Fatalf("querying foreign_keys pragma: %v", err)
	}
	if fkEnabled != 1 {
		t.Errorf("expected foreign_keys=1, got %d", fkEnabled)
	}
}
