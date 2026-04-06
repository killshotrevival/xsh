package host

import (
	"database/sql"
)

func Edit(db *sql.DB, hostID string) error {
	host, err := GetHostByName(db, hostID)
	if err != nil {
		return err
	}
	return createHost(db, host)
}
