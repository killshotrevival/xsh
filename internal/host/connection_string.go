package host

import (
	"database/sql"
	"fmt"
	"xsh/internal/tool"
)

func BuildConnectionString(identifier string, dbConnection *sql.DB) (string, error) {
	cHost, err := GetHostByName(dbConnection, identifier)
	if err != nil {
		return "", err
	}

	if cHost.ToolID == tool.SSHToolID {
		return buildSSHConnectionString(cHost, dbConnection)
	}

	return "", fmt.Errorf("TODO: Yet To implement")
}
