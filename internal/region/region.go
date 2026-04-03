package region

import (
	"database/sql"

	"github.com/google/uuid"
)

type Region struct {
	Id   uuid.UUID `json:"id"` //nolint:revive
	Name string    `json:"name"`
	// Tags []string  `json:"tags"`
}

var (
	updateRegionStmt = "UPDATE REGIONS SET NAME = ? WHERE ID = ?"
	insertRegionStmt = "INSERT INTO regions (id, name) VALUES (?, ?)"

	deleteRegionStmt       = "DELETE FROM REGIONS WHERE ID = ?"
	selectRegionStmt       = "SELECT ID, NAME FROM REGIONS"
	selectRegionByNameStmt = "SELECT ID, NAME FROM REGIONS WHERE NAME LIKE ?"
	getRegionIDByNameStmt  = "SELECT ID FROM REGIONS WHERE NAME = ?"
	GetRegionByIDStmt      = "SELECT ID, NAME FROM REGIONS WHERE ID = ?"

	getHostIDByRegionStmt = "SELECT ID FROM HOSTS WHERE REGION_ID = ?"
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

func (r *Region) Update(db *sql.DB) error {
	_, err := db.Exec(updateRegionStmt, r.Name, r.Id)
	return err
}

func (r *Region) Store(db *sql.DB) error {
	_, err := db.Exec(insertRegionStmt, r.Id, r.Name)
	return err
}
