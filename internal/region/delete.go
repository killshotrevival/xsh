package region

import (
	"database/sql"
	"fmt"

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
