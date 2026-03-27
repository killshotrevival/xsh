package identity

import (
	"database/sql"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
)

func PutIdentity(db *sql.DB, name, path string) error {
	// TODO: Check if the path provided contains a file or not
	// Make sure the path provided is abslute path not relative path
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
		log.Debugf("[identity] failed to create tag mapping for identity %q and tag %q: %v", identittyName, tagName, err)
		return err
	}

	return tm.Store(db)
}
