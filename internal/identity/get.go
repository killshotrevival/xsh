package identity

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/charmbracelet/log"
)

func PrintIdentities(db *sql.DB, identifier string) error {
	identities := []Identity{}
	var rows *sql.Rows
	var err error

	for _, placeholder := range []string{"name", "id", "path"} {
		if identifier == "*" {
			rows, err = db.Query("select id, name, path from identities")
		} else {
			rows, err = db.Query("select id, name, path from identities where "+placeholder+" like '%?%';", identifier)
		}
		if err != nil {
			log.Debugf("error occurred while fetching identities: %v", err)
			continue
		}

		for rows.Next() {
			id := Identity{}
			err := rows.Scan(&id.ID, &id.Name, &id.Path)
			if err != nil {
				log.Debugf("error occurred while reading identity: %v", err)
				continue
			}
			identities = append(identities, id)
		}
		if identifier == "*" {
			break
		}
	}
	log.Debug("Writing data to file")

	by, _ := json.Marshal(&identities)

	os.WriteFile("test.json", by, 0644)

	return nil
}
