package tag

import (
	"database/sql"

	"github.com/google/uuid"
)

type Tag struct {
	Id  uuid.UUID `json:"id"`
	Tag string    `json:"tag"`
}

type TagMapping struct {
	Id         uuid.UUID `json:"id"`
	TagId      uuid.UUID `json:"tag_id"`
	DataTypeId uuid.UUID `json:"data_type_id"`
}

var (
	insertTagStmt        = "insert into tags (id, tag) values (?, ?)"
	insertTagMappingStmt = "insert into tagmappings (id, tag_Id, data_type_id) values (?, ?, ?)"

	deleteTagStmt        = "delete from tags where id = ?"
	deleteTagMappingStmt = "delete from tagmappings where id = ?"

	getTagIdStmt      = "select id from tags where tag = ?"
	getTagStmt        = "select id, tag from tags where tag = ?"
	getTagMappingStmt = "select id, tag_id, data_type_id from tagmappings where tag_Id = ? and data_type_id = ?"
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

func NewTagMapping(tagId, dataTypeId uuid.UUID) (*TagMapping, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &TagMapping{
		Id:         id,
		TagId:      tagId,
		DataTypeId: dataTypeId,
	}, nil
}

func (tm *TagMapping) Store(db *sql.DB) error {
	_, err := db.Exec(insertTagMappingStmt, tm.Id, tm.TagId, tm.DataTypeId)
	return err
}

func (tm *TagMapping) Delete(db *sql.DB) error {
	_, err := db.Exec(deleteTagMappingStmt, tm.Id)
	return err
}
