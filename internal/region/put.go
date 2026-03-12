package region

import "database/sql"

func PutRegion(db *sql.DB, name string) error {
	region, err := NewRegion(name)
	if err != nil {
		return err
	}
	return region.Store(db)
}
