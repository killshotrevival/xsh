package host

import (
	"database/sql"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var (
	getHostIDByNameStmt     = "SELECT ID FROM HOSTS WHERE NAME = ?"
	getHostIDByAddressStmt  = "SELECT ID FROM HOSTS WHERE ADDRESS = ?"
	getHostByNameStmt       = "SELECT ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID FROM HOSTS WHERE NAME = ?"
	getHostByIDStmt         = "SELECT ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID FROM HOSTS WHERE ID = ?"
	getHostByJumphostIDStmt = "SELECT ID FROM HOSTS WHERE JUMPHOST_ID = ?"
	getJumphostName         = "SELECT NAME FROM HOSTS WHERE ID = ?"
	getShortHostStmt        = "SELECT ID, NAME FROM HOSTS"

	getHostStmt            = "SELECT H.ID, H.NAME, H.ADDRESS, H.PORT, H.USER, H.JUMPHOST_ID, H.REGION_ID, H.IDENTITY_ID, R.NAME AS REGION, I.PATH AS IDENTITYFILE FROM HOSTS AS H JOIN REGIONS AS R ON R.ID = H.REGION_ID JOIN IDENTITIES AS I ON I.ID = H.IDENTITY_ID"
	getHostWithNameStmt    = "SELECT H.ID, H.NAME, H.ADDRESS, H.PORT, H.USER, H.JUMPHOST_ID, H.REGION_ID, H.IDENTITY_ID, R.NAME AS REGION, I.PATH AS IDENTITYFILE FROM HOSTS AS H JOIN REGIONS AS R ON R.ID = H.REGION_ID JOIN IDENTITIES AS I ON I.ID = H.IDENTITY_ID WHERE h.NAME LIKE ?;"
	getHostWithAddressStmt = "SELECT H.ID, H.NAME, H.ADDRESS, H.PORT, H.USER, H.JUMPHOST_ID, H.REGION_ID, H.IDENTITY_ID, R.NAME AS REGION, I.PATH AS IDENTITYFILE FROM HOSTS AS H JOIN REGIONS AS R ON R.ID = H.REGION_ID JOIN IDENTITIES AS I ON I.ID = H.IDENTITY_ID WHERE h.ADDRESS LIKE ?;"
	getHostWithUserStmt    = "SELECT H.ID, H.NAME, H.ADDRESS, H.PORT, H.USER, H.JUMPHOST_ID, H.REGION_ID, H.IDENTITY_ID, R.NAME AS REGION, I.PATH AS IDENTITYFILE FROM HOSTS AS H JOIN REGIONS AS R ON R.ID = H.REGION_ID JOIN IDENTITIES AS I ON I.ID = H.IDENTITY_ID WHERE h.User LIKE ?;"
	deleteHostStmt         = "DELETE FROM HOSTS where ID = ?"
	insertHostStmt         = "INSERT INTO HOSTS (ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
)

type ShortHost struct {
	Id   uuid.UUID //nolint:revive
	Name string
}

type Host struct {
	Id         uuid.UUID     `json:"id"` //nolint:revive
	Name       string        `json:"name" comment:"Unique name of the host"`
	Address    string        `json:"address" comment:"Domain / IP address of the host without port"`
	Port       int           `json:"port" comment:"Port on which ssh connection will be created"`
	User       string        `json:"user" comment:"Remote user for creating the ssh connection"`
	RegionID   uuid.UUID     `json:"region_id" comment:"UUID of the region you want to connect this host to. You can find the id by printing the region table (xsg get 'r' '*')"`
	IdentityID uuid.UUID     `json:"identity_id" comment:"UUID of the Identity key you want to use for connecting with the host. You can find the id by printing the identity table (xsg get 'i' '*')"`
	JumphostID uuid.NullUUID `json:"jumphost_id" comment:"UUID of the host you want to use as jumphost. You can get the id by printing the host table (xsh get 'h' '*')"`
	// Tags         []string      `json:"tags"`
	Region       string `json:"region_name"`
	Jumphost     string `json:"jumphost_name"`
	IdentityFile string `json:"identitiy_file_name"`
}

func NewHost(name, address, user string, port int, regionID, identityID uuid.UUID, jumphostID uuid.NullUUID) (*Host, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	if port == 0 {
		log.Debug("[host] port not specified, defaulting to 22")
		port = 22
	}
	return &Host{
		Id:         id,
		Name:       name,
		Address:    address,
		Port:       port,
		User:       user,
		RegionID:   regionID,
		IdentityID: identityID,
		JumphostID: jumphostID,
	}, nil
}

func (h *Host) Store(db *sql.DB) error {
	if err := checkAddress(db, h.Address); err != nil {
		log.Warn("[host] a host with this address already exists, skipping insert")
		return nil
	}

	if err := checkName(db, h.Name); err != nil {
		log.Warn("[host] a host with this name already exists, skipping insert")
		return nil
	}
	_, err := db.Exec(insertHostStmt, h.Id, h.Name, h.Address, h.Port, h.User, h.RegionID, h.IdentityID, h.JumphostID)
	return err
}

func (h *Host) getJumphost(db *sql.DB) {
	jumpHostName := "-"

	if h.JumphostID.Valid {
		if err := db.QueryRow(getJumphostName, h.JumphostID).Scan(&jumpHostName); err != nil {
			if err == sql.ErrNoRows {
				jumpHostName = "No host present with ID attached"
			}
			jumpHostName = "DB error while checking"
		}
	}

	h.Jumphost = jumpHostName
}
