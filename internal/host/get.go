package host

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"xsh/internal/table"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var (
	noHostFoundError = "no host found with the given identifier"
)

func GetHostByName(db *sql.DB, identifier string) (*Host, error) {
	return getHost(
		db,
		getHostByNameStmt,
		identifier,
	)
}

func GetHostByID(db *sql.DB, identifier string) (*Host, error) {
	return getHost(
		db,
		GetHostByIDStmt,
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
		&host.RegionID,
		&host.IdentityID,
		&host.JumphostID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s : %v", noHostFoundError, err)
		}
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

func GetShortHosts(db *sql.DB) (*[]ShortHost, error) {
	rows, err := db.Query(getShortHostStmt)
	if err != nil {
		log.Debugf("error occurred while reading hosts from database: %v", err)
		return nil, err
	}

	var hosts []ShortHost
	for rows.Next() {
		var sh ShortHost

		if err := rows.Scan(&sh.Id, &sh.Name); err != nil {
			log.Debugf("error occurred while reading row: %v", err)
			return nil, err
		}

		hosts = append(hosts, sh)
	}

	return &hosts, nil
}

func Print(db *sql.DB, identifier string, outputFormat string) error {
	var rows *sql.Rows
	var err error

	idsAdded := []uuid.UUID{}
	data := [][]string{}
	printHost := []Host{}

	for _, placeholder := range []string{getHostWithNameStmt, getHostWithAddressStmt, getHostWithUserStmt} {
		if identifier == "*" {
			log.Debug("Printing all the hosts present in database")
			rows, err = db.Query(getHostStmt)
		} else {
			rows, err = db.Query(placeholder, "%"+identifier+"%")
		}
		if err != nil {
			log.Debugf("error occurred while fetching hosts: %v", err)
			continue
		}

		for rows.Next() {
			var (
				host Host
			)
			if err := rows.Scan(
				&host.Id,
				&host.Name,
				&host.Address,
				&host.Port,
				&host.User,
				&host.JumphostID,
				&host.RegionID,
				&host.IdentityID,
				&host.Region,
				&host.IdentityFile,
			); err != nil {
				log.Debugf("error occurred while reading host: %v", err)
				continue
			}

			if !slices.Contains(idsAdded, host.Id) {
				host.Tags, err = tag.GetTagsByDataTypeID(db, host.Id)
				if err != nil {
					host.Tags = []string{"error occurred while fetching"}
				}
				idsAdded = append(idsAdded, host.Id)
				host.getJumphost(db)

				printHost = append(printHost, host)
				data = append(data, []string{
					host.Name,
					fmt.Sprintf("%s:%d", host.Address, host.Port),
					host.Jumphost,
					host.User,
					host.Region,
					host.IdentityFile,
					tag.ToString(host.Tags),
				})

			}
		}
		if identifier == "*" {
			break
		}
	}

	switch strings.ToLower(outputFormat) {
	case "table":
		log.Debug("Printing data in table")
		t := table.NewTable(
			[]string{"NAME", "ADDRESS", "JUMPHOST", "USER", "REGION", "IDENTITY FILE", "TAGS"},
			data,
		)
		return t.Print()

	case "json":
		log.Debug("writing data in hosts.json file")
		by, _ := json.Marshal(&printHost)
		return os.WriteFile("hosts.json", by, 0644)
	default:
		return fmt.Errorf("invalid output format provided")
	}
}
