package region

import (
	"database/sql"

	"github.com/google/uuid"
)

type Region struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `josn:"name"`
}

var (
	CreateRegionTableStmt = `CREATE TABLE IF NOT EXISTS regions (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	slug TEXT NOT NULL
	)`

	insertRegionStmt = "INSERT INTO regions (id, name, slug) VALUES (?, ?, ?)"
)

func NewRegion(name, slug string) (*Region, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Region{
		Id:   id,
		Name: name,
		Slug: slug,
	}, nil
}

func (r *Region) Store(db *sql.DB) error {
	_, err := db.Exec(insertRegionStmt, r.Id, r.Name, r.Slug)
	return err
}
