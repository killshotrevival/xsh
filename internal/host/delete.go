package host

import (
	"database/sql"
	"fmt"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func Delete(db *sql.DB, identifier string) error {
	var h uuid.UUID

	err := db.QueryRow(getHostIDByNameStmt, identifier).Scan(&h)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no host found with the given name (%s)", identifier)
		}
		return err
	}

	if _, err = db.Exec(deleteHostStmt, h); err != nil {
		return err
	}

	// TODO: Delete tag mapping for the host

	return nil
}

func DeleteTagMapping(db *sql.DB, hostName, tagName string) error {
	host, nTag, err := getHostAndTag(db, hostName, tagName)
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
