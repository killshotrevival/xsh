package host

import (
	"database/sql"
)

func GetHostByName(db *sql.DB, identifier string) (*Host, error) {
	return getHost(
		db,
		getHostByNameStmt,
		identifier,
	)
}

func GetHostById(db *sql.DB, identifier string) (*Host, error) {
	return getHost(
		db,
		getHostByIdStmt,
		identifier,
	)
}

func getHost(db *sql.DB, queryString, identifier string) (*Host, error) {
	host := Host{}
	if err := db.QueryRow(queryString, identifier).Scan(
		&host.Id,
		&host.Name,
		&host.Address,
		&host.Port,
		&host.User,
		&host.RegionId,
		&host.IdentityId,
		&host.JumphostId,
	); err != nil {
		return nil, err
	}

	return &host, nil
}
