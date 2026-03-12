package region

import (
	"database/sql"
	"encoding/json"
	"os"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func PrintRegions(db *sql.DB, identifier string) error {
	var rows *sql.Rows
	var err error

	regions := []Region{}
	idsAdded := []uuid.UUID{}

	for _, placeholder := range []string{"name", "id", "slug"} {
		if identifier == "*" {
			log.Info("Printing all the regions present in database")
			rows, err = db.Query("select id, name, slug from Regions")
		} else {
			rows, err = db.Query("select id, name, slug from regions where "+placeholder+" like ?;", "%"+identifier+"%")
		}
		if err != nil {
			log.Debugf("error occurred while fetching Regions: %v", err)
			continue
		}

		for rows.Next() {
			r := Region{}
			err := rows.Scan(&r.Id, &r.Name, &r.Slug)
			if err != nil {
				log.Debugf("error occurred while reading identity: %v", err)
				continue
			}

			if !slices.Contains(idsAdded, r.Id) {
				idsAdded = append(idsAdded, r.Id)
				regions = append(regions, r)
			}
		}
		if identifier == "*" {
			break
		}
	}
	log.Debug("Writing data to file")

	by, _ := json.Marshal(&regions)

	os.WriteFile("region.json", by, 0644)

	return nil
}
