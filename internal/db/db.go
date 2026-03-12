package db

import (
	"database/sql"
	"os"
	"path/filepath"
	config "xsh/internal/config"
	"xsh/internal/identity"

	_ "github.com/mattn/go-sqlite3"
)

var (
	createRegionTableStmt = `CREATE TABLE IF NOT EXISTS regions (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	slug TEXT NOT NULL
	)`

	createHostTableStmt = `CREATE TABLE IF NOT EXISTS hosts (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT NOT NULL,
	user TEXT NOT NULL,
	region_id UUID NOT NULL,
	identity_id UUID NOT NULL,
	jumphost_id UUID
	)`
)

func GetDBPath() (string, error) {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "xsh.db"), nil
}

func CheckDB() (bool, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return false, err
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

func GetDB() (*sql.DB, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitDB() error {
	dbPath, err := GetDBPath()
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create tables or perform other initialization tasks here
	_, err = db.Exec(identity.CreateIdentityTableStmt)
	if err != nil {
		return err
	}

	_, err = db.Exec(createRegionTableStmt)
	if err != nil {
		return err
	}

	_, err = db.Exec(createHostTableStmt)
	if err != nil {
		return err
	}

	return nil
}
