package host

import (
	"database/sql"
	"fmt"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var (
	errJumphostDelete = fmt.Errorf("other hosts are using this resource as jumphost, can not proceed with deleting")
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

	rows, err := db.Query(getHostByJumphostIDStmt, h)
	if err != nil {
		log.Debugf("[host] error occurred while checking the jumphost mapping: %v", err)
		return err
	}

	defer rows.Close()

	if rows.Next() {
		return errJumphostDelete
	}

	if _, err = db.Exec(deleteHostStmt, h); err != nil {
		return err
	}

	return nil
}

func DeleteTagMapping(db *sql.DB, hostName, tagName string) error {
	host, nTag, err := getHostAndTag(db, hostName, tagName)
	if err != nil {
		return err
	}
	tm, err := tag.GetTagMapping(db, nTag.Id, host.Id)
	if err != nil {
		log.Debugf("[host] failed to retrieve tag mapping for host %q and tag %q: %v", hostName, tagName, err)
		return err
	}

	return tm.Delete(db)
}
