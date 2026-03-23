package region

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

func GetRegionByName(db *sql.DB, identifier string) (*Region, error) {
	region := Region{Name: identifier}

	if err := db.QueryRow(getRegionIDByNameStmt, identifier).Scan(&region.Id); err != nil {
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

func GetRegions(db *sql.DB) (*[]Region, error) {
	rows, err := db.Query(selectRegionStmt)
	if err != nil {
		return nil, err
	}

	var regions []Region
	for rows.Next() {
		var reg Region
		if err := rows.Scan(
			&reg.Id, &reg.Name,
		); err != nil {
			return nil, err
		}

		regions = append(regions, reg)
	}
	return &regions, nil
}

func Print(db *sql.DB, identifier, outputFormat string) error {
	var rows *sql.Rows
	var err error

	regions := []Region{}
	idsAdded := []uuid.UUID{}
	data := [][]string{}

	for _, stmt := range []string{selectRegionByNameStmt} {
		if identifier == "*" {
			log.Info("Printing all the regions present in database")
			rows, err = db.Query(selectRegionStmt)
		} else {
			rows, err = db.Query(stmt, "%"+identifier+"%")
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
				r.Tags, err = tag.GetTagsByDataTypeID(db, r.Id)
				if err != nil {
					r.Tags = []string{"error occurred while fetching"}
				}
				idsAdded = append(idsAdded, r.Id)
				regions = append(regions, r)
				data = append(data, []string{
					r.Id.String(),
					r.Name,
					tag.ToString(r.Tags),
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
		return table.NewTable(
			[]string{
				"ID",
				"NAME",
				"TAGS",
			},
			data,
		).Print()
	case "json":
		log.Debug("Writing data to file")
		by, _ := json.Marshal(&regions)
		return os.WriteFile("region.json", by, 0644)
	default:
		return fmt.Errorf("invalid output format received")
	}
}
