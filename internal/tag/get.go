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

var (
	noTagFoundError = "no tag found with given identifier"
)

func GetTag(db *sql.DB, identifier string) (*Tag, error) {
	tag := Tag{}
	if err := db.QueryRow(getTagWithTagStmt, identifier).Scan(&tag.Id, &tag.Tag); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: (%s)", noTagFoundError, identifier)
		}
		return nil, err
	}
	return &tag, nil
}

func GetTagsByDataTypeID(db *sql.DB, dataTypeID uuid.UUID) ([]string, error) {
	tags := []string{}

	rows, err := db.Query(getTagsByDataTypeIDStmt, dataTypeID)
	if err != nil {
		log.Debugf("[tag] failed to query tags for data type ID %q: %v", dataTypeID, err)
		return nil, err
	}

	for rows.Next() {
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			log.Debugf("[tag] failed to scan tag row from result set: %v", err)
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

// GetTagWithCreate will create a new with value equals to identifier if no tag is present in the database with the given name
func GetTagWithCreate(db *sql.DB, identifier string) (*Tag, error) {
	tag := &Tag{}
	if err := db.QueryRow(getTagWithTagStmt, identifier).Scan(&tag.Id, &tag.Tag); err != nil {
		if err == sql.ErrNoRows {
			log.Debug("[tag] no existing tag found for the given value, creating a new entry")
			tag, err = NewTag(identifier)
			if err != nil {
				log.Debugf("[tag] failed to create new tag with value %q", identifier)
				return nil, err
			}

			if err := tag.Store(db); err != nil {
				log.Debugf("[tag] failed to persist new tag %q to database: %v", identifier, err)
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
			log.Debug("[tag] listing all tags from the database")
			rows, err = db.Query(getTagStmt)
		} else {
			rows, err = db.Query(placeholder, "%"+identifier+"%")
		}
		if err != nil {
			log.Debugf("[tag] failed to query tags matching identifier %q: %v", identifier, err)
			continue
		}

		for rows.Next() {
			tag := Tag{}
			if err := rows.Scan(
				&tag.Id,
				&tag.Tag,
			); err != nil {
				log.Debugf("[tag] failed to scan tag row during listing: %v", err)
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
		log.Debug("[tag] exporting tag data to tags.json")
		by, _ := json.Marshal(&tags)
		return os.WriteFile("tags.json", by, 0600)
	default:
		return fmt.Errorf("invalid output format received")
	}
}
