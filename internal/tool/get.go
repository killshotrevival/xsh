package tool

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"xsh/internal/table"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func GetTools(db *sql.DB) (*[]Tool, error) {
	tools := []Tool{}

	rows, err := db.Query(getToolsStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := Tool{}
		if err := rows.Scan(&t.ID, &t.Name, &t.ConnectionString); err != nil {
			return nil, err
		}

		tools = append(tools, t)
	}

	return &tools, nil
}

func GetToolByID(db *sql.DB, identifier uuid.UUID) (*Tool, error) {
	t := Tool{}
	if err := db.QueryRow(getToolByIDStmt, identifier).Scan(&t.ID, &t.Name, &t.ConnectionString); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tool present with given identifier")
		}
		return nil, err
	}
	return &t, nil
}

func GetToolByName(db *sql.DB, identifier string) (*Tool, error) {
	t := Tool{}
	if err := db.QueryRow(getToolByNameStmt, identifier).Scan(&t.ID, &t.Name, &t.ConnectionString); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tool present with given identifier")
		}
		return nil, err
	}
	return &t, nil
}

func Print(db *sql.DB, identifier, outputFormat, outputFile string) error {
	var rows *sql.Rows
	var err error

	tools := []Tool{}
	idsAdded := []uuid.UUID{}

	data := [][]string{}

	for _, placeholder := range []string{getToolByNameStmt, getToolByIDStmt} {
		if identifier == "*" {
			rows, err = db.Query(getToolsStmt)
		} else {
			rows, err = db.Query(placeholder, "%"+identifier+"%")
		}
		if err != nil {
			log.Debugf("[tool] failed to query identities matching identifier %q: %v", identifier, err)
			continue
		}

		for rows.Next() {
			id := Tool{}
			err := rows.Scan(&id.ID, &id.Name, &id.ConnectionString)
			if err != nil {
				log.Debugf("[tool] failed to scan identity row during listing: %v", err)
				continue
			}
			if !slices.Contains(idsAdded, id.ID) {
				idsAdded = append(idsAdded, id.ID)
				tools = append(tools, id)
				data = append(data, []string{
					id.ID.String(),
					id.Name,
					id.ConnectionString,
				})
			}
		}
		if identifier == "*" {
			break
		}
	}
	switch strings.ToLower(outputFormat) {
	case "table":
		t := table.NewTable([]string{
			"ID", "NAME", "CONNECTION TEMPLATE",
		}, data)

		return t.Print()
	case "json":
		log.Debug("[tool] exporting identity data to json file")

		by, _ := json.Marshal(&tools)

		return os.WriteFile(outputFile, by, 0600)
	default:
		return fmt.Errorf("invalid output format provided")
	}
}
