package host

import (
	"database/sql"
	"fmt"
	"strings"
	"xsh/internal/tool"

	"github.com/charmbracelet/log"
)

func BuildConnectionString(dbConnection *sql.DB, identifier string, verbose bool) (string, error) {
	cHost, err := GetHostByName(dbConnection, identifier)
	if err != nil {
		return "", err
	}

	if cHost.ToolID == tool.SSHToolID {
		return buildSSHConnectionString(dbConnection, cHost, verbose)
	}

	return BuildFromTemplate(dbConnection, cHost, verbose)
}

func BuildFromTemplate(dbConnection *sql.DB, cHost *Host, verbose bool) (string, error) {
	if verbose {
		log.Warn("connection templates don't support verbose output generation right now")
	}

	to, err := tool.GetToolByID(dbConnection, cHost.ToolID)
	if err != nil {
		return "", err
	}

	connectionString := to.ConnectionString

	variables := extractVariablesFromTemplate(connectionString)

	for _, variable := range variables {
		value, err := cHost.GetValue(dbConnection, variable)
		if err != nil {
			return "", err
		}

		connectionString = strings.ReplaceAll(connectionString, fmt.Sprintf("${%s}", variable), value)
	}

	return connectionString, nil
}
