package identity

import (
	"database/sql"
	"os"
	"xsh/internal/utils"

	"charm.land/huh/v2"
)

func Edit(db *sql.DB, identityId string) error {
	id, err := GetIdentityByName(db, identityId)
	if err != nil {
		return err
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Identity Name").
				Description("Please select identity name").
				Value(&id.Name),

			huh.NewInput().
				Title("Identity Path").
				Description("Please select identity path").
				Value(&id.Path).Validate(func(s string) error {
				s, err := utils.ConvertToAbs(s)
				if err != nil {
					return err
				}
				_, err = os.Stat(s)
				return err

			}),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	return id.Update(db)
}
