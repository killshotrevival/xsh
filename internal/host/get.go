package host

import (
	"database/sql"
	"encoding/json"
	"os"
	"slices"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
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

func getHostAndTag(db *sql.DB, hostName, tagName string) (*Host, *tag.Tag, error) {
	host, err := GetHostByName(db, hostName)
	if err != nil {
		log.Debugf("error occurred while fetching host from given identifier(%s): %v", hostName, err)
		return nil, nil, err
	}

	nTag, err := tag.GetTagWithCreate(db, tagName)
	if err != nil {
		return nil, nil, err
	}

	return host, nTag, nil

}

func Print(db *sql.DB, identifier string) error {
	var rows *sql.Rows
	var err error

	hosts := []Host{}
	idsAdded := []uuid.UUID{}

	for _, placeholder := range []string{"name", "id"} {
		if identifier == "*" {
			log.Info("Printing all the hosts present in database")
			rows, err = db.Query("SELECT ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID FROM HOSTS")
		} else {
			rows, err = db.Query("SELECT ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID FROM HOSTS WHERE "+placeholder+" LIKE ?;", "%"+identifier+"%")
		}
		if err != nil {
			log.Debugf("error occurred while fetching hosts: %v", err)
			continue
		}

		for rows.Next() {
			host := Host{}
			if err := rows.Scan(
				&host.Id,
				&host.Name,
				&host.Address,
				&host.Port,
				&host.User,
				&host.RegionId,
				&host.IdentityId,
				&host.JumphostId,
			); err != nil {
				log.Debugf("error occurred while reading host: %v", err)
				continue
			}

			if !slices.Contains(idsAdded, host.Id) {
				idsAdded = append(idsAdded, host.Id)
				hosts = append(hosts, host)
			}
		}
		if identifier == "*" {
			break
		}
	}
	log.Debug("Writing data to file")

	by, _ := json.Marshal(&hosts)

	os.WriteFile("hosts.json", by, 0644)

	return nil
}
