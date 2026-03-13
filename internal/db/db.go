package db

import (
	"database/sql"
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	config "xsh/internal/config"

	"github.com/charmbracelet/log"

	_ "github.com/mattn/go-sqlite3"
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
		log.Debugf("error occurred while fetching the version number from the schema_version table: %v", err)
		return -1, err
	}
	return version, nil

}

func CheckAndApplyMigrations() error {
	newVersion := -1
	db, err := GetDB()
	if err != nil {
		log.Debugf("error occurred while connecting to database: %v", err)
		return err
	}

	defer db.Close()

	files, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		log.Debugf("error occurred while reading migrations directory: %v", err)
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

	log.Debug("checking which all migration files to pply based on version number")
	for _, file := range files {
		versionStr := strings.Split(file.Name(), "_")[0]

		version, err := strconv.Atoi(versionStr)
		if err != nil {
			log.Debugf("error occurred while converting migratiion file name to version number(%s): %v", file.Name(), err)
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
		log.Infof("Seems like there are %d number of database migrations that are not yet applied to the local database. Applying them ...", len(fileNames))
		if err := applyMigrations(db, fileNames); err != nil {
			return err
		}

		if _, err = db.Exec(updateSchemaVersionQuery, newVersion); err != nil {
			log.Debugf("error occurred while updating schema version to latest versiion %d -> %d", currentVersion, newVersion)
			return err
		}
	} else {
		log.Debug("No new migration file found for applying")
	}
	return nil
}

func applyMigrations(db *sql.DB, fileNames []string) error {
	for _, file := range fileNames {
		content, err := migrationFiles.ReadFile(file)
		if err != nil {
			log.Debugf("error occurred while reading migration file (%s): %v", file, err)
			return err
		}

		if _, err = db.Exec(string(content)); err != nil {
			log.Debugf("error occurred while applying migration (%s): %v", file, err)
			return err
		}
	}
	return nil
}

func InitDB() error {
	db, err := GetDB()
	if err != nil {
		log.Debugf("error occurred while connecting to database: %v", err)
		return err
	}

	defer db.Close()

	files, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		log.Debugf("error occurred while reading migrations directory: %v", err)
		return err
	}
	fileNames := []string{}

	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	log.Info("Applyting migrations")
	applyMigrations(db, fileNames)

	return nil
}
