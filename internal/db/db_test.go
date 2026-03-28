package db

import (
	"os"
	"path/filepath"
	"testing"
	"xsh/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetDBPath(t *testing.T) {

	path := utils.GetXSHTempDir(t)

	defer utils.RemoveTempDir(path, t)

	dbPath := filepath.Join(path, "xsh.db")

	if err := os.Setenv("XSH_DB_PATH", dbPath); err != nil {
		t.Fatalf("error occured while populating XSH_DB_PATH env: %v", err)
	}

	p, err := GetDBPath()
	if err != nil {
		t.Fatalf("error occurred while fetchinig DB PATh")
	}

	assert.Equal(t, dbPath, p)

}

func TestInitDB(t *testing.T) {

	path := utils.GetXSHTempDir(t)

	defer utils.RemoveTempDir(path, t)

	dbPath := filepath.Join(path, "xsh.db")

	if err := os.Setenv("XSH_DB_PATH", dbPath); err != nil {
		t.Fatalf("error occured while populating XSH_DB_PATH env: %v", err)
	}

	if err := InitDB(); err != nil {
		t.Fatalf("error occurred while initialising database: %v", err)
	}
}
