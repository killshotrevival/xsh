package identity

import (
	"database/sql"
	"os"
	"xsh/internal/tag"
	"xsh/internal/utils"

	"github.com/charmbracelet/log"
)

func PutIdentity(db *sql.DB, name, path string) error {

	path, err := utils.ConvertToAbs(path)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		log.Debugf("error occurred while trying to fetch identity file details: %v", err)
		return err
	}
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
