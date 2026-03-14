package tag

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func GetTag(db *sql.DB, identifier string) (*Tag, error) {
	tag := Tag{}
	if err := db.QueryRow(getTagStmt, identifier).Scan(&tag.Id, &tag.Tag); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tag found with given identifier: (%s)", identifier)
		}
		return nil, err
	}
	return &tag, nil
}

func GetTagsByDatatypeId(db *sql.DB, dataTypeId uuid.UUID) ([]string, error) {
	tags := []string{}

	rows, err := db.Query(getTagsByDatatypeIdStmt, dataTypeId)
	if err != nil {
		log.Debugf("error occurred while fetching tags for given datatype id(%s): %v", dataTypeId, err)
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

func GetTagMapping(db *sql.DB, tagId, datatypeId uuid.UUID) (*TagMapping, error) {
	tm := TagMapping{}
	if err := db.QueryRow(getTagMappingStmt, tagId, datatypeId).Scan(&tm.Id, &tm.TagId, &tm.DataTypeId); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tag mapping found with given identifiers: (%s | %s)", tagId, datatypeId)
		}
		return nil, err
	}
	return &tm, nil
}

func GetTagWithCreate(db *sql.DB, identifier string) (*Tag, error) {
	tag := &Tag{}
	if err := db.QueryRow(getTagStmt, identifier).Scan(&tag.Id, &tag.Tag); err != nil {
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

func Print(db *sql.DB, identifier string) error {
	var rows *sql.Rows
	var err error

	tags := []Tag{}
	idsAdded := []uuid.UUID{}

	for _, placeholder := range []string{"name", "id"} {
		if identifier == "*" {
			log.Info("Printing all the tags present in database")
			rows, err = db.Query("SELECT ID, TAG FROM TAGS")
		} else {
			rows, err = db.Query("SELECT ID, TAG FROM TAGS WHERE "+placeholder+" LIKE ?;", "%"+identifier+"%")
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
				tags = append(tags, tag)
			}
		}
		if identifier == "*" {
			break
		}
	}
	log.Debug("Writing data to file")

	by, _ := json.Marshal(&tags)

	os.WriteFile("tags.json", by, 0644)

	return nil
}
