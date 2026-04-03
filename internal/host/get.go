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

func checkAddress(db *sql.DB, address string) error {
	var hID string
	if address == "" {
		return fmt.Errorf("empty string address received")
	}
	if err := db.QueryRow(getHostIDByAddressStmt, address).Scan(&hID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return fmt.Errorf("host present in database with address provided")
}

func checkName(db *sql.DB, name string) error {
	var hID string
	if name == "" {
		return fmt.Errorf("empty string name received")
	}
	if err := db.QueryRow(getHostIDByNameStmt, name).Scan(&hID); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return fmt.Errorf("host present in database with name provided")
}

func GetHostByID(db *sql.DB, identifier string) (*Host, error) {
	return getHost(
		db,
		getHostByIDStmt,
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
		&host.ExtraFlags,
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
		log.Debugf("[host] failed to retrieve host by identifier %q: %v", hostName, err)
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
		log.Debugf("[host] failed to query short host list from database: %v", err)
		return nil, err
	}

	var hosts []ShortHost
	for rows.Next() {
		var sh ShortHost

		if err := rows.Scan(&sh.Id, &sh.Name); err != nil {
			log.Debugf("[host] failed to scan short host row from result set: %v", err)
			return nil, err
		}

		hosts = append(hosts, sh)
	}

	return &hosts, nil
}

func Print(db *sql.DB, identifier, outputFormat, outputFile string) error {
	var rows *sql.Rows
	var err error

	idsAdded := []uuid.UUID{}
	data := [][]string{}
	printHost := []Host{}

	for _, placeholder := range []string{getHostWithNameStmt, getHostWithAddressStmt, getHostWithUserStmt} {
		if identifier == "*" {
			log.Debug("[host] listing all hosts from the database")
			rows, err = db.Query(getHostStmt)
		} else {
			rows, err = db.Query(placeholder, "%"+identifier+"%")
		}
		if err != nil {
			log.Debugf("[host] failed to query hosts matching identifier %q: %v", identifier, err)
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
				&host.ExtraFlags,
				&host.RegionID,
				&host.IdentityID,
				&host.Region,
				&host.IdentityFile,
			); err != nil {
				log.Debugf("[host] failed to scan host row during listing: %v", err)
				continue
			}

			if !slices.Contains(idsAdded, host.Id) {
				// TODO: Freezed until further development
				// host.Tags, err = tag.GetTagsByDataTypeID(db, host.Id)
				// if err != nil {
				// 	host.Tags = []string{"error occurred while fetching"}
				// }
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
					host.ExtraFlags,
					// tag.ToString(host.Tags),
				})

			}
		}
		if identifier == "*" {
			break
		}
	}

	switch strings.ToLower(outputFormat) {
	case "table":
		log.Debug("[host] rendering host data as table")
		t := table.NewTable(
			[]string{"NAME", "ADDRESS", "JUMPHOST", "USER", "REGION", "IDENTITY FILE", "EXTRA FLAGS"}, // "TAGS"

			data,
		)
		return t.Print()

	case "json":
		log.Debug("[host] exporting host data to json file", "outputfile", outputFile)
		by, _ := json.Marshal(&printHost)
		return os.WriteFile(outputFile, by, 0600)
	default:
		return fmt.Errorf("invalid output format provided")
	}
}
