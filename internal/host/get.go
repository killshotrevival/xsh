package host

import (
	"database/sql"
	"fmt"
)

func GetHostByName(db *sql.DB, identifier string) (*Host, error) {
	return getHost(
		db,
		getHostStmt,
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
	rows, err := db.Query(queryString, identifier)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("no host present with the given identifier (%s)", identifier)
	}

	if err = rows.Scan(
		&host.Id,
		&host.Name,
		&host.Address,
		&host.User,
		&host.RegionId,
		&host.IdentityId,
		&host.JumphostId,
	); err != nil {
		return nil, err
	}

	return &host, nil
}
