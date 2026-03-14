package identity

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"xsh/internal/tag"

	"github.com/google/uuid"

	"github.com/charmbracelet/log"
)

func GetIdentityById(db *sql.DB, identifier uuid.UUID) (*Identity, error) {
	id := Identity{}

	if err := db.QueryRow(getIdentityByIdStmt, identifier).Scan(&id.Id, &id.Name, &id.Path); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no identity found with the given identifier (%s)", identifier)
		}
		return nil, err
	}

	return &id, nil
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

func Print(db *sql.DB, identifier string) error {
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
				id.Tags, err = tag.GetTagsByDatatypeId(db, id.Id)
				if err != nil {
					id.Tags = []string{"error occurred while fetching"}
				}
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
