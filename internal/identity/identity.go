package identity

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"xsh/internal/utils"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var (
	insertIdentityStmt = "INSERT INTO IDENTITIES (ID, NAME, PATH) VALUES (?, ?, ?)"

	deleteIdentityStmt = "DELETE FROM identities WHERE ID = ?"

	getIdentityIDByPathStmt = "SELECT ID FROM IDENTITIES WHERE PATH = ?"
	getIdentityStmt         = "SELECT ID, NAME, PATH FROM IDENTITIES"
	getIdentityIDByNameStmt = "SELECT ID FROM IDENTITIES WHERE NAME = ?"
	getIdentityByNameStmt   = "SELECT ID, NAME, PATH FROM IDENTITIES WHERE NAME = ?"
	getIdentityByPathStmt   = "SELECT ID, NAME, PATH FROM IDENTITIES WHERE PATH = ?"
	getIdentityByIDStmt     = "SELECT ID, NAME, PATH FROM IDENTITIES WHERE ID = ?"

	getHostIDByIdentityStmt = "SELECT ID FROM HOSTS WHERE IDENTITY_ID = ?"

	sshKeyMarkers = []string{
		"-----BEGIN OPENSSH PRIVATE KEY-----",
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN EC PRIVATE KEY-----",
		"-----BEGIN DSA PRIVATE KEY-----",
	}
)

type Identity struct {
	Id   uuid.UUID `json:"id"` //nolint:revive
	Name string    `json:"name" comment:"Name of the identity file"`
	Path string    `json:"path" comment:"Absolute path of the identity file"`
	// Tags []string  `json:"tags"`
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

func CheckOrCreateIdentity(path string, db *sql.DB) (uuid.UUID, error) {
	var id uuid.UUID
	path, err := utils.ConvertToAbs(path)
	if err != nil {
		return uuid.UUID{}, err
	}
	if err := db.QueryRow(getIdentityIDByPathStmt, path).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			log.Debugf("No path exists with path (%s), creating a new one", path)

			if _, err := os.Stat(path); err != nil {
				log.Debugf("[identity] error occurred while trying to fetch identity file details: %v", err)
				return uuid.UUID{}, err
			}

			i, err := NewIdentity(filepath.Base(path), path)
			if err != nil {
				log.Debugf("[identity] error occurred while trying to create a new identity: %v", err)
				return uuid.UUID{}, err
			}

			if err := i.Store(db); err != nil {
				log.Debugf("[identity] error occurred while trying to save the identity to database: %v", err)
				return uuid.UUID{}, err
			}

			return i.Id, nil
		}
	}
	return id, nil
}

func (i *Identity) Store(db *sql.DB) error {
	status, err := i.ExistsInDb(db)
	if err != nil {
		return nil
	}

	if status {
		log.Warn("[identity] identity with this path already exists in the database, skipping insert")
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

	log.Debugf("[identity] discovered SSH directory at: %s", sshHomeDir)

	if err := filepath.WalkDir(sshHomeDir, func(path string, d os.DirEntry, err error) error {
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

			log.Debugf("[identity] persisting SSH identity file to database: %s", identity.Path)
			err = identity.Store(db)
			if err != nil {
				log.Errorf("Error storing identity for %s: %v", path, err)
				return nil
			}

		}

		return nil
	}); err != nil {
		log.Debugf("[identity] failed to walk SSH home directory %q: %v", sshHomeDir, err)
		return err
	}

	return nil
}
