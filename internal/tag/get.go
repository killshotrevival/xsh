package tag

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"xsh/internal/table"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func GetTag(db *sql.DB, identifier string) (*Tag, error) {
	tag := Tag{}
	if err := db.QueryRow(getTagWithTagStmt, identifier).Scan(&tag.Id, &tag.Tag); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tag found with given identifier: (%s)", identifier)
		}
		return nil, err
	}
	return &tag, nil
}

func GetTagsByDataTypeID(db *sql.DB, dataTypeID uuid.UUID) ([]string, error) {
	tags := []string{}

	rows, err := db.Query(getTagsByDataTypeIDStmt, dataTypeID)
	if err != nil {
		log.Debugf("error occurred while fetching tags for given datatype id(%s): %v", dataTypeID, err)
		return nil, err
	}

	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			log.Debugf("error occurred while reading tag from database: %v", err)
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func GetTagMapping(db *sql.DB, tagID, dataTypeID uuid.UUID) (*Mapping, error) {
	tm := Mapping{}
	if err := db.QueryRow(getTagMappingStmt, tagID, dataTypeID).Scan(&tm.Id, &tm.TagID, &tm.DataTypeID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tag mapping found with given identifiers: (%s | %s)", tagID, dataTypeID)
		}
		return nil, err
	}
	return &tm, nil
}

func GetTagWithCreate(db *sql.DB, identifier string) (*Tag, error) {
	tag := &Tag{}
	if err := db.QueryRow(getTagWithTagStmt, identifier).Scan(&tag.Id, &tag.Tag); err != nil {
		if err == sql.ErrNoRows {
			log.Debug("no tag exists in database with given value, creating a new one")
			tag, err = NewTag(identifier)
			if err != nil {
				log.Debugf("error occurred while creating new tag")
				return nil, err
			}

			if err := tag.Store(db); err != nil {
				log.Debugf("error occurred while storing new tag to database: %v", err)
				return nil, err
			}
		}
		return nil, err
	}
	return tag, nil
}

func Print(db *sql.DB, identifier string, outputFormat string) error {
	var rows *sql.Rows
	var err error

	tags := []Tag{}
	idsAdded := []uuid.UUID{}
	data := [][]string{}

	for _, placeholder := range []string{getTagWithTagStmt} {
		if identifier == "*" {
			log.Debug("Printing all the tags present in database")
			rows, err = db.Query(getTagStmt)
		} else {
			rows, err = db.Query(placeholder, "%"+identifier+"%")
		}
		if err != nil {
			log.Debugf("error occurred while fetching tags: %v", err)
			continue
		}

		for rows.Next() {
			tag := Tag{}
			if err := rows.Scan(
				&tag.Id,
				&tag.Tag,
			); err != nil {
				log.Debugf("error occurred while reading tag: %v", err)
				continue
			}

			if !slices.Contains(idsAdded, tag.Id) {
				idsAdded = append(idsAdded, tag.Id)
				data = append(data, []string{tag.Tag})
				tags = append(tags, tag)
			}
		}
		if identifier == "*" {
			break
		}
	}
	switch strings.ToLower(outputFormat) {
	case "table":
		return table.NewTable(
			[]string{"TAGS"},
			data,
		).Print()
	case "json":
		log.Debug("Writing data to file")
		by, _ := json.Marshal(&tags)
		return os.WriteFile("tags.json", by, 0644)
	default:
		return fmt.Errorf("invalid output format received")
	}
}
