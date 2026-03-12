package identity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/google/uuid"

	"github.com/charmbracelet/log"
)

func GetIdentityById(db *sql.DB, identifier uuid.UUID) (*Identity, error) {
	id := Identity{}

	rows, err := db.Query(getIdentityByIdStmt, identifier)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("no identity present with the given identifier (%s)", identifier)
	}

	if err = rows.Scan(&id.Id, &id.Name, &id.Path); err != nil {
		return nil, err
	}

	return &id, nil
}

func PrintIdentities(db *sql.DB, identifier string) error {
	var rows *sql.Rows
	var err error

	identities := []Identity{}
	idsAdded := []uuid.UUID{}

	for _, placeholder := range []string{"name", "id", "path"} {
		if identifier == "*" {
			rows, err = db.Query("select id, name, path from identities")
		} else {
			rows, err = db.Query("select id, name, path from identities where "+placeholder+" like ?;", "%"+identifier+"%")
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
				idsAdded = append(idsAdded, id.Id)
				identities = append(identities, id)
			}
		}
		if identifier == "*" {
			break
		}
	}
	log.Debug("Writing data to file")

	by, _ := json.Marshal(&identities)

	os.WriteFile("identity.json", by, 0644)

	return nil
}
