package tag

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Tag struct {
	Id  uuid.UUID `json:"id"` //nolint:revive
	Tag string    `json:"tag"`
}

type Mapping struct {
	Id         uuid.UUID `json:"id"` //nolint:revive
	TagID      uuid.UUID `json:"tag_id"`
	DataTypeID uuid.UUID `json:"data_type_id"`
}

var (
	insertTagStmt        = "INSERT INTO TAGS (ID, TAG) VALUES (?, ?)"
	insertTagMappingStmt = "INSERT INTO TAGMAPPINGS (ID, TAG_ID, DATA_TYPE_ID) VALUES (?, ?, ?)"

	deleteTagStmt        = "DELETE FROM TAGS WHERE ID = ?"
	deleteTagMappingStmt = "DELETE FROM TAGMAPPINGS WHERE ID = ?"

	getTagStmt              = "SELECT ID, TAG FROM TAGS"
	getTagIDStmt            = "SELECT ID FROM TAGS WHERE TAG = ?"
	getTagWithTagStmt       = "SELECT ID, TAG FROM TAGS WHERE TAG = ?"
	getTagsByDataTypeIDStmt = "SELECT T.TAG FROM TAGS AS T JOIN TAGMAPPINGS AS TM ON T.ID = TM.TAG_ID WHERE TM.DATA_TYPE_ID = ?"
	getTagMappingStmt       = "SELECT ID, TAG_ID, DATA_TYPE_ID FROM TAGMAPPINGS WHERE TAG_ID = ? AND DATA_TYPE_ID = ?"
)

func NewTag(tag string) (*Tag, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Tag{
		Id:  id,
		Tag: tag,
	}, nil
}

func (t *Tag) Store(db *sql.DB) error {
	_, err := db.Exec(insertTagStmt, t.Id, t.Tag)
	return err
}

func NewTagMapping(tagID, dataTypeID uuid.UUID) (*Mapping, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Mapping{
		Id:         id,
		TagID:      tagID,
		DataTypeID: dataTypeID,
	}, nil
}

func (tm *Mapping) Store(db *sql.DB) error {
	_, err := db.Exec(insertTagMappingStmt, tm.Id, tm.TagID, tm.DataTypeID)
	return err
}

func (tm *Mapping) Delete(db *sql.DB) error {
	_, err := db.Exec(deleteTagMappingStmt, tm.Id)
	return err
}

func ToString(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	finalStr := tags[0]

	for _, item := range tags[1:] {
		finalStr = fmt.Sprintf("%s, %s", finalStr, item)
	}
	return finalStr
}
