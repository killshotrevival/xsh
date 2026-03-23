package tag

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func Delete(db *sql.DB, identifier string) error {
	var id uuid.UUID

	if err := db.QueryRow(getTagIDStmt, identifier).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no tag found with given identifier (%s)", identifier)
		}
		return err
	}

	if _, err := db.Exec(deleteTagStmt, id); err != nil {
		return err
	}
	return nil
}
