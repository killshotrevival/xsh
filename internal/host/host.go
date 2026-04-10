package host

import (
	"database/sql"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"xsh/internal/identity"
	"xsh/internal/region"
	"xsh/internal/tool"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

var (
	getHostIDByNameStmt     = "SELECT ID FROM HOSTS WHERE NAME = ?"
	getHostIDByAddressStmt  = "SELECT ID FROM HOSTS WHERE ADDRESS = ?"
	getHostByNameStmt       = "SELECT ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID, TOOL_ID, EXTRA_FLAGS FROM HOSTS WHERE NAME = ?"
	getHostByIDStmt         = "SELECT ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID, TOOL_ID, EXTRA_FLAGS FROM HOSTS WHERE ID = ?"
	getHostByJumphostIDStmt = "SELECT ID FROM HOSTS WHERE JUMPHOST_ID = ?"
	getJumphostName         = "SELECT NAME FROM HOSTS WHERE ID = ?"
	getShortHostStmt        = "SELECT ID, NAME FROM HOSTS"
	getHostStmt             = "SELECT H.ID, H.NAME, H.ADDRESS, H.PORT, H.USER, H.JUMPHOST_ID, H.TOOL_ID, H.EXTRA_FLAGS, H.REGION_ID, H.IDENTITY_ID, R.NAME AS REGION, I.PATH AS IDENTITYFILE FROM HOSTS AS H JOIN REGIONS AS R ON R.ID = H.REGION_ID JOIN IDENTITIES AS I ON I.ID = H.IDENTITY_ID"
	getHostWithNameStmt     = "SELECT H.ID, H.NAME, H.ADDRESS, H.PORT, H.USER, H.JUMPHOST_ID, H.TOOL_ID, H.EXTRA_FLAGS, H.REGION_ID, H.IDENTITY_ID, R.NAME AS REGION, I.PATH AS IDENTITYFILE FROM HOSTS AS H JOIN REGIONS AS R ON R.ID = H.REGION_ID JOIN IDENTITIES AS I ON I.ID = H.IDENTITY_ID WHERE h.NAME LIKE ?;"
	getHostWithAddressStmt  = "SELECT H.ID, H.NAME, H.ADDRESS, H.PORT, H.USER, H.JUMPHOST_ID, H.TOOL_ID, H.EXTRA_FLAGS, H.REGION_ID, H.IDENTITY_ID, R.NAME AS REGION, I.PATH AS IDENTITYFILE FROM HOSTS AS H JOIN REGIONS AS R ON R.ID = H.REGION_ID JOIN IDENTITIES AS I ON I.ID = H.IDENTITY_ID WHERE h.ADDRESS LIKE ?;"
	getHostWithUserStmt     = "SELECT H.ID, H.NAME, H.ADDRESS, H.PORT, H.USER, H.JUMPHOST_ID, H.TOOL_ID, H.EXTRA_FLAGS, H.REGION_ID, H.IDENTITY_ID, R.NAME AS REGION, I.PATH AS IDENTITYFILE FROM HOSTS AS H JOIN REGIONS AS R ON R.ID = H.REGION_ID JOIN IDENTITIES AS I ON I.ID = H.IDENTITY_ID WHERE h.User LIKE ?;"
	deleteHostStmt          = "DELETE FROM HOSTS where ID = ?"

	updateHostStmt = "UPDATE HOSTS SET NAME = ?, ADDRESS = ?, PORT = ?, USER = ?, REGION_ID = ?, IDENTITY_ID = ?, JUMPHOST_ID = ?, TOOL_ID = ?, EXTRA_FLAGS = ? WHERE ID = ?"
	insertHostStmt = "INSERT INTO HOSTS (ID, NAME, ADDRESS, PORT, USER, REGION_ID, IDENTITY_ID, JUMPHOST_ID, TOOL_ID, EXTRA_FLAGS) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
)

type ShortHost struct {
	Id   uuid.UUID //nolint:revive
	Name string
}

type Host struct {
	Id           uuid.UUID     `json:"id"` //nolint:revive
	Name         string        `json:"name" comment:"Unique name of the host"`
	Address      string        `json:"address" comment:"Domain / IP address of the host without port"`
	Port         int           `json:"port" comment:"Port on which ssh connection will be created"`
	User         string        `json:"user" comment:"Remote user for creating the ssh connection"`
	RegionID     uuid.UUID     `json:"region_id" comment:"UUID of the region you want to connect this host to. You can find the id by printing the region table (xsg get 'r' '*')"`
	IdentityID   uuid.UUID     `json:"identity_id" comment:"UUID of the Identity key you want to use for connecting with the host. You can find the id by printing the identity table (xsg get 'i' '*')"`
	JumphostID   uuid.NullUUID `json:"jumphost_id" comment:"UUID of the host you want to use as jumphost. You can get the id by printing the host table (xsh get 'h' '*')"`
	ToolID       uuid.UUID     `json:"tool_id" comment:"UUID of the tool you want to use for making connection."`
	ExtraFlags   string        `json:"extra_flags" comment:"Extra ssh flgs, except XSH internal ones"`
	Region       string        `json:"region_name"`
	Jumphost     string        `json:"jumphost_name"`
	IdentityFile string        `json:"identitiy_file_name"`
}

func (h *Host) GetValue(db *sql.DB, field string) (string, error) {
	switch field {
	case "address":
		return h.Address, nil
	case "port":
		return strconv.Itoa(h.Port), nil
	case "user":
		return h.User, nil
	case "extra_flags":
		return h.ExtraFlags, nil
	case "identitiy_file_path":
		id, err := identity.GetIdentityByID(db, h.IdentityID)
		if err != nil {
			return "", fmt.Errorf("error occurred while fetching identity path: %v", err)
		}
		return id.Path, nil
	default:
		return "", fmt.Errorf("%s is not a valid place holder found in the template", field)

	}
}

// Template will have variables in the format of
// ${variable name here}, this function will extract all such variables from the string
// and returns them in a list
func extractVariablesFromTemplate(template string) []string {
	variables := []string{}

	re := regexp.MustCompile(`\$\{([^}]+)\}`)
	// FindAllStringSubmatch returns a slice of slices:
	// matches[n][0] is the full string (e.g., "${port}")
	// matches[n][1] is the first capturing group (e.g., "port")
	matches := re.FindAllStringSubmatch(template, -1)

	for _, match := range matches {
		if len(match) > 1 {
			if !slices.Contains(variables, match[1]) {
				variables = append(variables, match[1])
			}
		}
	}

	return variables
}

// Function is used for validating if the template connection string contains all the valid variables or not.
func ValidateConnectionString(connectionTemplate string) error {
	validKeys := []string{
		"address",
		"port",
		"user",
		"extra_flags",
		"identitiy_file_path",
	}
	variables := extractVariablesFromTemplate(connectionTemplate)

	for _, variable := range variables {
		if !slices.Contains(validKeys, variable) {
			return fmt.Errorf("%s is not a valid place holder found in the template", variable)
		}
	}

	return nil
}

// Create a host with default values filled up
func GetDefaultHost() (*Host, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Host{
		Id:         id,
		Port:       22,
		User:       "root",
		RegionID:   region.DefaultregionID,
		IdentityID: identity.DefaultIdentityID,
		ToolID:     tool.SSHToolID,
	}, nil

}

// Create a new host with the provided values
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

func (h *Host) Update(db *sql.DB) error {
	_, err := db.Exec(updateHostStmt, h.Name, h.Address, h.Port, h.User, h.RegionID, h.IdentityID, h.JumphostID, h.ToolID, h.ExtraFlags, h.Id)
	return err
}

func (h *Host) Store(db *sql.DB) error {
	if err := checkAddress(db, h.Address); err != nil {
		log.Warn("[host] a host with this address already exists, skipping insert", "address", h.Address)
		return nil
	}

	if err := checkName(db, h.Name); err != nil {
		log.Warn("[host] a host with this name already exists, skipping insert", "name", h.Name)
		return nil
	}

	if h.IdentityID == [16]byte{} {
		log.Warn("[host] invalid identity assigned, skipping insert", "identityId", h.IdentityID)
		return nil
	}

	if h.RegionID == [16]byte{} {
		log.Warn("[host] invalid region assigned, skipping insert", "regionId", h.RegionID)
		return nil
	}

	_, err := db.Exec(insertHostStmt, h.Id, h.Name, h.Address, h.Port, h.User, h.RegionID, h.IdentityID, h.JumphostID, h.ToolID, h.ExtraFlags)
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

func (h *Host) UpdateExtraFlags(flag string) {

	newFlags := fmt.Sprintf("%s %s", h.ExtraFlags, flag)

	if err := validateExtraFlags(newFlags); err != nil {
		log.Debug("Flag cannot be added in the extra flags", "error", err)
		return
	}

	h.ExtraFlags = newFlags

}

func (h *Host) UpdateUserAddress(s string) {
	if strings.Contains(s, "@") {
		addressSplit := strings.Split(s, "@")
		h.User = addressSplit[0]
		h.Address = addressSplit[1]
	}
}

// This function is used for validating the flags string should not contain any
// of the restricted flags
func validateExtraFlags(flags string) error {
	if len(flags) == 0 {
		// No extra flags present to validate
		return nil
	}

	for _, flag := range []string{"-J", "-v", "-p", "-i", "-l", "-q", "-V"} {
		if strings.Contains(flags, flag) {
			return fmt.Errorf("cannot use `%s` in extra flags, as this flag is handled by XSH internally", flag)
		}
	}

	if strings.Contains(flags, "-4") && strings.Contains(flags, "-6") {
		return fmt.Errorf("-4 and -6 flags cannot be used together")
	}

	if strings.Contains(flags, "-A") && strings.Contains(flags, "-a") {
		return fmt.Errorf("-a and -A flags cannot be used together")
	}

	return nil
}
