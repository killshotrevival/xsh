package region

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"xsh/internal/tag"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func GetRegionByName(db *sql.DB, identifier string) (*Region, error) {
	region := Region{Name: identifier}

	if err := db.QueryRow(getRegionIdByNameStmt, identifier).Scan(&region.Id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no region present with given identifier (%s): %v", identifier, err)
		}
		return nil, err
	}

	return &region, nil
}

func getRegionAndTag(db *sql.DB, identittyName, tagName string) (*Region, *tag.Tag, error) {
	host, err := GetRegionByName(db, identittyName)
	if err != nil {
		log.Debugf("error occurred while fetching region from given identifier(%s): %v", identittyName, err)
		return nil, nil, err
	}

	nTag, err := tag.GetTagWithCreate(db, tagName)
	if err != nil {
		return nil, nil, err
	}

	return host, nTag, nil

}

func PrintRegions(db *sql.DB, identifier string) error {
	var rows *sql.Rows
	var err error

	regions := []Region{}
	idsAdded := []uuid.UUID{}

	for _, placeholder := range []string{"name", "id"} {
		if identifier == "*" {
			log.Info("Printing all the regions present in database")
			rows, err = db.Query("select id, name from Regions")
		} else {
			rows, err = db.Query("select id, name from regions where "+placeholder+" like ?;", "%"+identifier+"%")
		}
		if err != nil {
			log.Debugf("error occurred while fetching Regions: %v", err)
			continue
		}

		for rows.Next() {
			r := Region{}
			err := rows.Scan(&r.Id, &r.Name)
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
