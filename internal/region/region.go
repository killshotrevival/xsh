package region

import (
	"database/sql"

	"github.com/google/uuid"
)

type Region struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Tags []string  `json:"tags"`
}

var (
	insertRegionStmt = "INSERT INTO regions (id, name) VALUES (?, ?)"

	deleteRegionStmt      = "DELETE REGIONS WHERE ID = ?"
	getRegionIdByNameStmt = "SELECT ID FROM REGIONS WHERE NAME = ?"
)

func NewRegion(name string) (*Region, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Region{
		Id:   id,
		Name: name,
	}, nil
}

func (r *Region) Store(db *sql.DB) error {
	_, err := db.Exec(insertRegionStmt, r.Id, r.Name)
	return err
}
