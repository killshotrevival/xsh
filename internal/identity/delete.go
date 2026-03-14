package identity

import (
	"database/sql"
	"fmt"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func Delete(db *sql.DB, identifier string) error {
	var i uuid.UUID

	if err := db.QueryRow(getIdentityIdByNameStmt, identifier).Scan(&i); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no identity found with the given identifier (%s)", identifier)
		}
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
		log.Debugf("error occurred while fetching tag mapping object: %v", err)
		return err
	}

	return tm.Delete(db)
}
