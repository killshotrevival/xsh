package host

import (
	"database/sql"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var (
	getHostIdByNameStmt = "SELECT ID FROM HOSTS WHERE NAME = ?"
	getHostByNameStmt   = "SELECT ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID FROM HOSTS WHERE NAME = ?"
	getHostByIdStmt     = "SELECT ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID FROM HOSTS WHERE ID = ?"
	getJumphostName     = "SELECT NAME FROM HOSTS WHERE ID = ?"

	printHostStmt  = "SELECT H.ID, H.NAME, H.ADDRESS, H.PORT, H.USER, H.JUMPHOST_ID, H.REGION_ID, H.IDENTITY_ID, R.NAME AS REGION, I.PATH AS IDENTITYFILE FROM HOSTS AS H JOIN REGIONS AS R ON R.ID = H.REGION_ID JOIN IDENTITIES AS I ON I.ID = H.IDENTITY_ID"
	deleteHostStmt = "DELETE FROM HOSTS where ID = ?"
	insertHostStmt = "INSERT INTO HOSTS (ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
)

type Host struct {
	Id           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	Address      string        `json:"address"`
	Port         int           `json:"port"`
	User         string        `json:"user"`
	RegionId     uuid.UUID     `json:"region_id"`
	IdentityId   uuid.UUID     `json:"identity_id"`
	JumphostId   uuid.NullUUID `json:"jumphost_id"`
	Tags         []string      `json:"tags"`
	Region       string        `json:"region_name"`
	Jumphost     string        `json:"jumphost_name"`
	IdentityFile string        `json:"identitiy_file_name"`
}

func NewHost(name, address, user string, port int, region_id, identityId uuid.UUID, jumphostId uuid.NullUUID) (*Host, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	if port == 0 {
		log.Debug("Port number found is 0, defaulting it to 22")
		port = 22
	}
	return &Host{
		Id:         id,
		Name:       name,
		Address:    address,
		Port:       port,
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

	_, err = db.Exec(insertHostStmt, h.Id, h.Name, h.Address, h.Port, h.User, h.RegionId, h.IdentityId, h.JumphostId)
	return err
}

func (h *Host) tagsString() string {
	if len(h.Tags) == 0 {
		return ""
	}
	finalStr := h.Tags[0]

	for _, item := range h.Tags[1:] {
		finalStr = fmt.Sprintf("%s, %s", finalStr, item)
	}
	return finalStr
}

func (h *Host) getJumphost(db *sql.DB) {
	jumpHostName := "-"

	if h.JumphostId.Valid {
		if err := db.QueryRow(getJumphostName, h.JumphostId).Scan(&jumpHostName); err != nil {
			if err == sql.ErrNoRows {
				jumpHostName = "No host present with ID attached"
			}
			jumpHostName = "DB error while checking"
		}
	}

	h.Jumphost = jumpHostName
}
