package cmd

import (
	"database/sql"
	"fmt"
	"slices"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"
	"xsh/internal/theme"

	"charm.land/huh/v2"
	"github.com/charmbracelet/log"
)

func selectResource(dbConnection *sql.DB) (string, []string, error) {
	var (
		resource string
		idLists  []string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&resource).
				Height(5).
				Title("Select Resource").
				Description("Please select the resource to delete").
				OptionsFunc(func() []huh.Option[string] {
					opts := []huh.Option[string]{}

					for key := range resourceDeleteMapping {
						opts = append(opts, huh.NewOption(key, key))
					}
					return opts
				}, nil),
		),
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				TitleFunc(func() string {
					return fmt.Sprintf("Select %s(s)", resource)
				}, nil).
				DescriptionFunc(func() string {
					return fmt.Sprintf("Please select all the %s(s) to delete", resource)
				}, nil).OptionsFunc(
				func() []huh.Option[string] {
					opts := []huh.Option[string]{}

					switch resource {
					case "host":
						allHosts, err := host.GetShortHosts(dbConnection)
						if err != nil {
							log.Debugf("[delete] failed to retrieve host list for interactive selection: %v", err)
							return []huh.Option[string]{
								huh.NewOption("error occurred while trying to select hosts", "-1"),
							}
						}
						for _, host := range *allHosts {
							opts = append(opts, huh.NewOption(host.Name, host.Name))
						}
					case "region":
						allRegions, err := region.GetRegions(dbConnection)
						if err != nil {
							log.Debugf("[delete] failed to retrieve region list for interactive selection: %v", err)
							return []huh.Option[string]{
								huh.NewOption("error occurred while trying to select regions", "-1"),
							}
						}
						for _, region := range *allRegions {
							opts = append(opts, huh.NewOption(region.Name, region.Name))
						}
					case "identity":
						allIDs, err := identity.GetIdentity(dbConnection)
						if err != nil {
							log.Debugf("[delete] failed to retrieve identity list for interactive selection: %v", err)
							return []huh.Option[string]{
								huh.NewOption("error occurred while trying to select identities", "-1"),
							}
						}
						for _, id := range *allIDs {
							opts = append(opts, huh.NewOption(id.Name, id.Name))
						}

					}
					return opts

				}, nil).
				Value(&idLists).Validate(func(s []string) error {
				if slices.Contains(s, "-1") {
					return fmt.Errorf("invalid value selected")
				}
				return nil
			}),
		),
	).WithTheme(huh.ThemeFunc(theme.XSH))

	if err := form.Run(); err != nil {
		return "", nil, err
	}

	if len(idLists) == 0 {
		return "", nil, fmt.Errorf("No reosurce seleceted")
	}

	return resource, idLists, nil
}
