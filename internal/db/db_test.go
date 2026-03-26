package db

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func removeTempDir(path string, t *testing.T) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Error removing temp directory %s: %v\n", path, err)
	}
}

func TestGetDBPath(t *testing.T) {
	path, err := os.MkdirTemp("", "xsh_temp_dir_*")
	if err != nil {
		t.Fatalf("error ocurred while creating temp dir for testing: %v", err)
	}

	defer removeTempDir(path, t)

	dbPath := filepath.Join(path, "xsh.db")

	if err = os.Setenv("XSH_DB_PATH", dbPath); err != nil {
		t.Fatalf("error occured while populating XSH_DB_PATH env: %v", err)
	}

	p, err := GetDBPath()
	if err != nil {
		t.Fatalf("error occurred while fetchinig DB PATh")
	}

	assert.Equal(t, p, dbPath)

}

func TestInitDB(t *testing.T) {
	path, err := os.MkdirTemp("", "xsh_temp_dir_*")
	if err != nil {
		t.Fatalf("error ocurred while creating temp dir for testing: %v", err)
	}

	defer removeTempDir(path, t)

	dbPath := filepath.Join(path, "xsh.db")

	if err = os.Setenv("XSH_DB_PATH", dbPath); err != nil {
		t.Fatalf("error occured while populating XSH_DB_PATH env: %v", err)
	}

	if err := InitDB(); err != nil {
		t.Fatalf("error occurred while initialising database: %v", err)
	}
}
