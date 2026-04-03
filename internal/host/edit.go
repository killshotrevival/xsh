package host

import (
	"database/sql"
)

func Edit(db *sql.DB, hostId string) error {
	host, err := GetHostByName(db, hostId)
	if err != nil {
		return err
	}
	return createHost(db, host)
}
