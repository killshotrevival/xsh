package db

import (
	"database/sql"
	"os"
	"path/filepath"
	config "xsh/internal/config"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"

	_ "github.com/mattn/go-sqlite3"
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

	_, err = db.Exec(region.CreateRegionTableStmt)
	if err != nil {
		return err
	}

	_, err = db.Exec(host.CreateHostTableStmt)
	if err != nil {
		return err
	}

	return nil
}
