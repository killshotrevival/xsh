package tool

import (
	"database/sql"

	"github.com/google/uuid"
)

var (
	SSHToolID = uuid.MustParse("0a398c97-4525-4416-903d-3662b8de8850")

	getToolsStmt      = "SELECT ID, NAME, CONNECTION_STRING FROM TOOLS"
	getToolByIDStmt   = "SELECT ID, NAME, CONNECTION_STRING FROM tools WHERE ID = ?"
	getToolByNameStmt = "SELECT ID, NAME, CONNECTION_STRING FROM tools WHERE NAME = ?"

	deleteToolStmt = "DELETE FROM TOOLS WHERE ID = ?"

	insertToolStmt = "INSERT INTO TOOLS (ID, NAME, CONNECTION_STRING) VALUES (?, ?, ?)"

	getHostIDByToolStmt = "SELECT ID FROM HOSTS WHERE TOOL_ID = ?"
)

type Tool struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	ConnectionString string    `json:"connection_string"`
}

func NewTool(name, connectionString string) (*Tool, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Tool{
		ID:               id,
		Name:             name,
		ConnectionString: connectionString,
	}, nil
}

func (t *Tool) Store(db *sql.DB) error {
	_, err := db.Exec(insertToolStmt, t.ID, t.Name, t.ConnectionString)
	return err
}
