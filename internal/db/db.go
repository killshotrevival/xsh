package db

import (
	"database/sql"
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	config "xsh/internal/config"

	"github.com/charmbracelet/log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver registration
)

var (

	//go:embed migrations/*.sql
	migrationFiles embed.FS

	tableExistsQuery = `
	SELECT name
	FROM sqlite_master
	WHERE type='table' AND name=?;
	`

	migrationVersionCheckQuery = `select version from schema_version;`
	updateSchemaVersionQuery   = `update schema_version set version = ?`
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

func GetTestDB(t *testing.T) *sql.DB {
	// ":memory:" creates a fresh DB in RAM for every call
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	files, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		t.Fatalf("error occurred while reading migrations directory: %v", err)
	}
	fileNames := []string{}

	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	if err := applyMigrations(db, fileNames, false); err != nil {
		t.Fatalf("error occurred while applying migrations: %v", err)
	}

	return db
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

func createBackup() error {

	dbPath, err := GetDBPath()
	if err != nil {
		return err
	}

	// Create destination path: same dir + ".bck"
	dst := dbPath + ".bck"

	// Open source file
	in, err := os.Open(dbPath)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create destination file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy contents
	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	// Ensure data is flushed
	if err := out.Sync(); err != nil {
		return err
	}

	return nil
}

func TableExists(db *sql.DB, tableName string) (bool, error) {
	var name string
	err := db.QueryRow(tableExistsQuery, tableName).Scan(&name)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func getCurrentAppliedMigrationVersion(db *sql.DB) (int, error) {
	var version int
	if err := db.QueryRow(migrationVersionCheckQuery).Scan(&version); err != nil {
		log.Debugf("[database] failed to retrieve current migration version from schema_version table: %v", err)
		return -1, err
	}
	return version, nil

}

func CheckAndApplyMigrations() error {
	newVersion := -1
	db, err := GetDB()
	if err != nil {
		log.Debugf("[database] failed to establish database connection for migration check: %v", err)
		return err
	}

	defer db.Close()

	files, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		log.Debugf("[database] failed to read embedded migrations directory: %v", err)
		return nil
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	currentVersion, err := getCurrentAppliedMigrationVersion(db)
	if err != nil {
		return err
	}
	fileNames := []string{}

	log.Debug("[database] evaluating pending migrations against current schema version")
	for _, file := range files {
		versionStr := strings.Split(file.Name(), "_")[0]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			log.Debugf("[database] failed to parse version number from migration filename %q: %v", file.Name(), err)
			return err
		}

		if version > currentVersion {
			if version > newVersion {
				newVersion = version
			}
			fileNames = append(fileNames, file.Name())
		}

	}

	if len(fileNames) > 0 {
		log.Debugf("[database] found %d pending migration(s) to apply to the local database", len(fileNames))
		if err := applyMigrations(db, fileNames, true); err != nil {
			return err
		}
		return updateSchemaVersion(db, newVersion)
	}
	log.Debug("[database] schema is up to date, no pending migrations to apply")
	return nil
}

func updateSchemaVersion(db *sql.DB, version int) error {
	if _, err := db.Exec(updateSchemaVersionQuery, version); err != nil {
		log.Debugf("[database] failed to update schema version to %d: %v", version, err)
		return err
	}
	return nil
}

func applyMigrations(db *sql.DB, fileNames []string, backup bool) error {
	if backup {
		log.Debug("[database] creating a backup file before applying migrations")
		if err := createBackup(); err != nil {
			log.Warnf("error occurred while creating a backup file for the database")
			return err
		}
	}
	for _, file := range fileNames {
		log.Debugf("[database] applying Migrations from %s file", file)
		content, err := migrationFiles.ReadFile("migrations/" + file)
		if err != nil {
			log.Debugf("[database] failed to read migration file %q: %v", file, err)
			return err
		}

		if _, err = db.Exec(string(content)); err != nil {
			log.Debugf("[database] failed to execute migration %q: %v", file, err)
			return err
		}
	}
	return nil
}

func InitDB() error {
	db, err := GetDB()
	if err != nil {
		log.Debugf("[database] failed to establish database connection during initialization: %v", err)
		return err
	}

	defer db.Close()

	files, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		log.Debugf("[database] failed to read embedded migrations directory during initialization: %v", err)
		return err
	}
	fileNames := []string{}
	latestVersion := 0

	for _, file := range files {
		fileNames = append(fileNames, file.Name())

		versionStr := strings.Split(file.Name(), "_")[0]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			log.Debugf("[database] failed to parse version number from migration filename %q: %v", file.Name(), err)
			return err
		}

		if version > latestVersion {
			log.Debugf("[database] Updating latest version to: %d", version)
			latestVersion = version
		}
	}

	log.Debugf("[database] applying initial database migrations")
	if err := applyMigrations(db, fileNames, false); err != nil {
		log.Debugf("[database] failed to apply initial migrations: %v", err)
		return err
	}

	return updateSchemaVersion(db, latestVersion)

}
