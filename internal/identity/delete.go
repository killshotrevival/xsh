package identity

import (
	"database/sql"
	"fmt"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func checkIdentity(db *sql.DB, identityID string) error {
	var hID string
	if err := db.QueryRow(getHostIDByIdentityStmt, identityID).Scan(&hID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return fmt.Errorf("host present in the database with identity provided")
}

func Delete(db *sql.DB, identifier string) error {
	var i uuid.UUID

	if err := db.QueryRow(getIdentityIDByNameStmt, identifier).Scan(&i); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no identity found with the given identifier (%s)", identifier)
		}
		return err
	}

	if err := checkIdentity(db, i.String()); err != nil {
		return err
	}

	if _, err := db.Exec(deleteIdentityStmt, i); err != nil {
		return err
	}
	return nil
}

func DeleteTagMapping(db *sql.DB, identittyName, tagName string) error {
	host, nTag, err := getIdentityAndTag(db, identittyName, tagName)
	if err != nil {
		return err
	}
	tm, err := tag.GetTagMapping(db, nTag.Id, host.Id)
	if err != nil {
		log.Debugf("[identity] failed to retrieve tag mapping for identity %q and tag %q: %v", identittyName, tagName, err)
		return err
	}

	return tm.Delete(db)
}
