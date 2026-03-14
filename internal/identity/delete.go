package identity

import (
	"database/sql"
	"fmt"

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
