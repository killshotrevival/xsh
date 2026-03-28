package identity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"xsh/internal/table"
	"xsh/internal/tag"

	"github.com/google/uuid"

	"github.com/charmbracelet/log"
)

var (
	noIdentityFoundErr = "no identity found with the given identifier"
)

func GetIdentityByID(db *sql.DB, identifier uuid.UUID) (*Identity, error) {
	id := Identity{}

	if err := db.QueryRow(GetIdentityByIDStmt, identifier).Scan(&id.Id, &id.Name, &id.Path); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no identity found with the given identifier (%s)", identifier)
		}
		return nil, err
	}

	return &id, nil
}

func GetIdentity(db *sql.DB) (*[]Identity, error) {
	rows, err := db.Query(getIdentityStmt)
	if err != nil {
		log.Debugf("[identity] failed to query identity list from database: %v", err)
		return nil, err
	}

	var ids []Identity
	for rows.Next() {
		var id Identity
		if err := rows.Scan(&id.Id, &id.Name, &id.Path); err != nil {
			log.Debugf("[identity] failed to scan identity row from result set: %v", err)
			return nil, err
		}
		ids = append(ids, id)
	}
	return &ids, nil
}

func GetIdentityByName(db *sql.DB, identifier string) (*Identity, error) {
	id := Identity{}

	if err := db.QueryRow(getIdentityByNameStmt, identifier).Scan(&id.Id, &id.Name, &id.Path); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s (%s)", noIdentityFoundErr, identifier)
		}
		return nil, err
	}

	return &id, nil
}

func getIdentityAndTag(db *sql.DB, identittyName, tagName string) (*Identity, *tag.Tag, error) {
	host, err := GetIdentityByName(db, identittyName)
	if err != nil {
		log.Debugf("[identity] failed to retrieve identity by name %q: %v", identittyName, err)
		return nil, nil, err
	}

	nTag, err := tag.GetTagWithCreate(db, tagName)
	if err != nil {
		return nil, nil, err
	}

	return host, nTag, nil

}

func Print(db *sql.DB, identifier, outputFormat, outputFile string) error {
	var rows *sql.Rows
	var err error

	identities := []Identity{}
	idsAdded := []uuid.UUID{}

	data := [][]string{}

	for _, placeholder := range []string{getIdentityByNameStmt, getIdentityByPathStmt} {
		if identifier == "*" {
			rows, err = db.Query(getIdentityStmt)
		} else {
			rows, err = db.Query(placeholder, "%"+identifier+"%")
		}
		if err != nil {
			log.Debugf("[identity] failed to query identities matching identifier %q: %v", identifier, err)
			continue
		}

		for rows.Next() {
			id := Identity{}
			err := rows.Scan(&id.Id, &id.Name, &id.Path)
			if err != nil {
				log.Debugf("[identity] failed to scan identity row during listing: %v", err)
				continue
			}
			if !slices.Contains(idsAdded, id.Id) {
				// TODO: Freezed until further development
				// id.Tags, err = tag.GetTagsByDataTypeID(db, id.Id)
				// if err != nil {
				// 	id.Tags = []string{"error occurred while fetching"}
				// }
				idsAdded = append(idsAdded, id.Id)
				identities = append(identities, id)
				data = append(data, []string{
					id.Id.String(),
					id.Name,
					id.Path,
					// tag.ToString(id.Tags),
				})
			}
		}
		if identifier == "*" {
			break
		}
	}
	switch strings.ToLower(outputFormat) {
	case "table":
		t := table.NewTable([]string{
			"ID", "NAME", "PATH",
		}, data)

		return t.Print()
	case "json":
		log.Debug("[identity] exporting identity data to json file")

		by, _ := json.Marshal(&identities)

		return os.WriteFile(outputFile, by, 0644)
	default:
		return fmt.Errorf("invalid output format provided")
	}
}
