package host

import (
	"database/sql"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var (
	CreateHostTableStmt = `CREATE TABLE IF NOT EXISTS hosts (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT NOT NULL,
	user TEXT NOT NULL,
	region_id UUID NOT NULL,
	identity_id UUID NOT NULL,
	jumphost_id UUID
	)`

	getHostStmt     = "select id, name, address, user, region_id, identity_id, jumphost_id from hosts where name = ?"
	getHostByIdStmt = "select id, name, address, user, region_id, identity_id, jumphost_id from hosts where id = ?"

	insertHostStmt = "INSERT INTO hosts (id, name, address, user, region_id, identity_id, jumphost_id) VALUES (?, ?, ?, ?, ?, ?, ?)"
)

type Host struct {
	Id         uuid.UUID     `json:"id"`
	Name       string        `json:"name"`
	Address    string        `json:"address"`
	User       string        `json:"user"`
	RegionId   uuid.UUID     `json:"region_id"`
	IdentityId uuid.UUID     `json:"identity_id"`
	JumphostId uuid.NullUUID `json:"jumphost_id"`
}

func NewHost(name, address, user string, region_id, identityId uuid.UUID, jumphostId uuid.NullUUID) (*Host, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &Host{
		Id:         id,
		Name:       name,
		Address:    address,
		User:       user,
		IdentityId: identityId,
		JumphostId: jumphostId,
	}, nil
}

func (h *Host) Store(db *sql.DB) error {
	rows, err := db.Query("select id from hosts where address = ?", h.Address)
	if err != nil {
		return err
	}
	if rows.Next() {
		log.Debug("Host with this address already exists")
		return nil
	}

	_, err = db.Exec(insertHostStmt, h.Id, h.Name, h.Address, h.User, h.RegionId, h.IdentityId, h.JumphostId)
	return err
}
