package tool

import (
	"database/sql"
	"fmt"
)

func Delete(db *sql.DB, identifier string) error {
	rows, err := db.Query(getHostIDByToolStmt, identifier)
	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		return fmt.Errorf("hosts present connected with this tool. Cannot proceed with deleting")
	}
	_, err = db.Exec(deleteToolStmt, identifier)
	return err
}
