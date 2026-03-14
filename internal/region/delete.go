package region

import (
	"database/sql"
	"fmt"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func Delete(db *sql.DB, identifier string) error {
	var id uuid.UUID
	if err := db.QueryRow(getRegionIdByNameStmt, identifier).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no region found with given name")
		}
		return err
	}

	if _, err := db.Exec(deleteRegionStmt, id); err != nil {
		return err
	}
	return nil
}

func DeleteTagMapping(db *sql.DB, identittyName, tagName string) error {
	host, nTag, err := getRegionAndTag(db, identittyName, tagName)
	if err != nil {
		return err
	}
	tm, err := tag.NewTagMapping(nTag.Id, host.Id)
	if err != nil {
		log.Debugf("error occurred while creating new tag mapping object; %v", err)
		return err
	}

	return tm.Delete(db)
}
