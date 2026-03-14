package identity

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var (
	CreateIdentityTableStmt = `CREATE TABLE IF NOT EXISTS identities (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	path TEXT NOT NULL
	)`

	insertIdentityStmt = "INSERT INTO IDENTITIES (ID, NAME, PATH) VALUES (?, ?, ?)"

	deleteIdentityStmt = "DELETE FROM identities WHERE ID = ?"

	getIdentityIdByNameStmt = "SELECT ID FROM IDENTITIES WHERE NAME = ?"
	getIdentityByNameStmt   = "SELECT ID, NAME, PATH FROM IDENTITIES WHERE NAME = ?"
	getIdentityByPathStmt   = "SELECT ID, NAME, PATH FROM IDENTITIES WHERE PATH = ?"
	getIdentityByIdStmt     = "SELECT ID, NAME, PATH FROM IDENTITIES WHERE ID = ?"

	sshKeyMarkers = []string{
		"-----BEGIN OPENSSH PRIVATE KEY-----",
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN EC PRIVATE KEY-----",
		"-----BEGIN DSA PRIVATE KEY-----",
	}
)

type Identity struct {
	Id   uuid.UUID `json:"id"`
	Name string    `josn:"name"`
	Path string    `json:"path"`
}

func NewIdentity(name, path string) (*Identity, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Identity{
		Id:   id,
		Name: name,
		Path: path,
	}, nil
}

func (i *Identity) Store(db *sql.DB) error {
	status, err := i.ExistsInDb(db)
	if err != nil {
		return nil
	}

	if status {
		log.Info("Identity with path already exists in the table")
		return nil
	}
	_, err = db.Exec(insertIdentityStmt, i.Id, i.Name, i.Path)
	return err
}

func (i *Identity) ExistsInDb(db *sql.DB) (bool, error) {
	rows, err := db.Query(getIdentityByPathStmt, i.Path)
	if err != nil {
		return false, fmt.Errorf("error occurred while checking identity execits in database: %w", err)
	}

	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func containsSSHKey(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read first few lines only (keys start at top)
	for i := 0; i < 10; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		for _, marker := range sshKeyMarkers {
			if strings.Contains(line, marker) {
				return true
			}
		}
	}

	return false
}

func InitIdentityStore(db *sql.DB) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	_ = homeDir // Use homeDir as needed

	sshHomeDir := filepath.Join(homeDir, ".ssh")

	log.Infof("SSH Dir found: %s", sshHomeDir)

	filepath.WalkDir(sshHomeDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.Type().IsRegular() {
			return nil
		}

		if containsSSHKey(path) {
			identity, err := NewIdentity(filepath.Base(path), path)
			if err != nil {
				log.Errorf("Error creating identity for %s: %v", path, err)
				return nil
			}

			log.Infof("Storing Identity file to database: %s", identity.Path)
			err = identity.Store(db)
			if err != nil {
				log.Errorf("Error storing identity for %s: %v", path, err)
				return nil
			}

		}

		return nil
	})

	return nil
}
