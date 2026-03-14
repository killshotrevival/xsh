package host

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func Delete(db *sql.DB, identifier string) error {
	var h uuid.UUID

	err := db.QueryRow(getHostIdByNameStmt, identifier).Scan(&h)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no host found with the given name (%s)", identifier)
		}
		return err
	}

	if _, err = db.Exec(deleteHostStmt, h); err != nil {
		return err
	}

	return nil
}
