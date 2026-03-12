package region

import "database/sql"

func PutRegion(db *sql.DB, name, slug string) error {
	region, err := NewRegion(name, slug)
	if err != nil {
		return err
	}
	return region.Store(db)
}
