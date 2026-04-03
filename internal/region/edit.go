package region

import (
	"database/sql"

	"charm.land/huh/v2"
)

func Edit(db *sql.DB, regionID string) error {
	reg, err := GetRegionByName(db, regionID)
	if err != nil {
		return err
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Region Name").
				Description("Please enter region name").
				Value(&reg.Name),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	return reg.Update(db)
}
