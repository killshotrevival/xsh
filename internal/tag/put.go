package tag

import (
	"database/sql"
)

func Put(db *sql.DB, tag string) error {
	nTag, err := NewTag(tag)
	if err != nil {
		return err
	}
	return nTag.Store(db)
}
