package identity

import (
	"database/sql"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
)

func PutIdentity(db *sql.DB, name, path string) error {
	id, err := NewIdentity(name, path)
	if err != nil {
		return nil
	}
	return id.Store(db)
}

func PutTagMapping(db *sql.DB, identittyName, tagName string) error {
	host, nTag, err := getIdentityAndTag(db, identittyName, tagName)
	if err != nil {
		return err
	}
	tm, err := tag.NewTagMapping(nTag.Id, host.Id)
	if err != nil {
		log.Debugf("error occurred while creating new tag mapping object; %v", err)
		return err
	}

	return tm.Store(db)
}
