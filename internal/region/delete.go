package region

import (
	"database/sql"
	"fmt"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func checkHost(db *sql.DB, regionID string) error {
	var hID string
	if err := db.QueryRow(getHostIDByRegionStmt, regionID).Scan(&hID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return fmt.Errorf("host present in the database attached with provided region")
}

func Delete(db *sql.DB, identifier string) error {
	var id uuid.UUID
	if err := db.QueryRow(getRegionIDByNameStmt, identifier).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no region found with given name")
		}
		return err
	}

	if err := checkHost(db, id.String()); err != nil {
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
	tm, err := tag.GetTagMapping(db, nTag.Id, host.Id)
	if err != nil {
		log.Debugf("[region] failed to retrieve tag mapping for region %q and tag %q: %v", identittyName, tagName, err)
		return err
	}

	return tm.Delete(db)
}
