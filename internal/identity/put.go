package identity

import "database/sql"

func PutIdentity(db *sql.DB, name, path string) error {
	id, err := NewIdentity(name, path)
	if err != nil {
		return nil
	}
	return id.Store(db)
}
