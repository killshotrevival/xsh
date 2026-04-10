package tool

import "database/sql"

func PutTool(db *sql.DB, name, connectionString string) error {
	t, err := NewTool(name, connectionString)
	if err != nil {
		return err
	}
	return t.Store(db)
}
