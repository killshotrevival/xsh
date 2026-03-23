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
		log.Debugf("error occurred while reading identity from database: %v", err)
		return nil, err
	}

	var ids []Identity
	for rows.Next() {
		var id Identity
		if err := rows.Scan(&id.Id, &id.Name, &id.Path); err != nil {
			log.Debugf("error occurred while reading identity row: %v", err)
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
			return nil, fmt.Errorf("no identity found with the given identifier (%s)", identifier)
		}
		return nil, err
	}

	return &id, nil
}

func getIdentityAndTag(db *sql.DB, identittyName, tagName string) (*Identity, *tag.Tag, error) {
	host, err := GetIdentityByName(db, identittyName)
	if err != nil {
		log.Debugf("error occurred while fetching identity from given identifier(%s): %v", identittyName, err)
		return nil, nil, err
	}

	nTag, err := tag.GetTagWithCreate(db, tagName)
	if err != nil {
		return nil, nil, err
	}

	return host, nTag, nil

}

func Print(db *sql.DB, identifier string, outputFormat string) error {
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
			log.Debugf("error occurred while fetching identities: %v", err)
			continue
		}

		for rows.Next() {
			id := Identity{}
			err := rows.Scan(&id.Id, &id.Name, &id.Path)
			if err != nil {
				log.Debugf("error occurred while reading identity: %v", err)
				continue
			}
			if !slices.Contains(idsAdded, id.Id) {
				id.Tags, err = tag.GetTagsByDataTypeID(db, id.Id)
				if err != nil {
					id.Tags = []string{"error occurred while fetching"}
				}
				idsAdded = append(idsAdded, id.Id)
				identities = append(identities, id)
				data = append(data, []string{
					id.Id.String(),
					id.Name,
					id.Path,
					tag.ToString(id.Tags),
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
			"ID", "NAME", "PATH", "TAGS",
		}, data)

		return t.Print()
	case "json":
		log.Debug("Writing data to file")

		by, _ := json.Marshal(&identities)

		return os.WriteFile("identity.json", by, 0644)
	default:
		return fmt.Errorf("invalid output format provided")
	}
}
