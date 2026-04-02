package host

import (
	"database/sql"
	"fmt"
	"strconv"
	"xsh/internal/identity"
	"xsh/internal/region"
	"xsh/internal/theme"

	"charm.land/huh/v2"
	"github.com/google/uuid"
)

func InteractivePut(db *sql.DB) error {
	var createOption int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("Choose how you want to create a new host").
				Description("Please choose how you want to proceed with creating a new host").
				Options(
					huh.NewOption("Clone from an existing host", 1),
					huh.NewOption("Create a new host from scratch", 2),
				).Value(&createOption),
		),
	).WithTheme(huh.ThemeFunc(theme.XSH))

	if err := form.Run(); err != nil {
		return err
	}

	if createOption == 1 {
		return cloneHost(db)
	}

	return createHost(db)
}

func cloneHost(db *sql.DB) error {
	var (
		cloneHostID string
	)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Host").
				Description("Please select the host you want to clone from").
				OptionsFunc(func() []huh.Option[string] {
					hosts, err := GetShortHosts(db)
					if err != nil {
						return []huh.Option[string]{
							huh.NewOption(
								fmt.Sprintf("error occurred while fetchinig hosts: %v", err), "-1",
							),
						}
					}

					if len(*hosts) == 0 {
						return []huh.Option[string]{
							huh.NewOption("no host present to select", "0"),
						}
					}

					opts := []huh.Option[string]{}
					for _, sh := range *hosts {
						opts = append(opts, huh.NewOption(
							sh.Name, sh.Id.String(),
						))
					}
					return opts

				}, nil).Value(&cloneHostID).Validate(func(s string) error {
				if s == "-1" {
					return fmt.Errorf("error occurred while selecting host, please exit and retry")
				}
				return nil
			}),
		),
	).WithTheme(huh.ThemeFunc(theme.XSH))

	if err := form.Run(); err != nil {
		return err
	}

	if cloneHostID == "0" {
		return createHost(db)
	}

	host, err := GetHostByID(db, cloneHostID)
	if err != nil {
		return err
	}

	host.Name = ""

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Host Name").
				Description("Please enter a unique name for your host that is easy to remember").
				Value(&host.Name).Validate(func(s string) error {
				return checkName(db, s)
			}),

			huh.NewInput().
				Title("Host Address").
				Description("Please enter the hostname / IP address of the host to connect").
				Value(&host.Address).Validate(func(s string) error {
				return checkAddress(db, s)
			}),
		),
	).WithTheme(huh.ThemeFunc(theme.XSH))

	if err := form.Run(); err != nil {
		return err
	}

	host.Id, err = uuid.NewRandom()
	if err != nil {
		return err
	}

	return host.Store(db)
}

func createHost(db *sql.DB) error {
	var (
		host             Host
		regionIDString   string
		identityIDString string
		jumphostIDString string
		err              error
	)

	portString := "22"
	host.User = "root"
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Host Name").
				Description("Please enter a unique name for your host that is easy to remember").
				Value(&host.Name).Validate(func(s string) error {
				return checkName(db, s)
			}),

			huh.NewInput().
				Title("Host Address").
				Description("Please enter the hostname / IP address of the host to connect").
				Value(&host.Address).Validate(func(s string) error {
				return checkAddress(db, s)
			}),

			huh.NewInput().
				Title("Host Port").
				Description("Please enter the port on which host accepts the connection").
				Value(&portString).Placeholder(portString),

			huh.NewInput().
				Title("Remote User").
				Description("Please enter the remote username").
				Value(&host.User).Placeholder(host.User),

			huh.NewInput().
				Title("SSH Flags").
				Description("Extra SSH flags, that you would like to use for this host").
				Value(&host.ExtraFlags).Validate(validateExtraFlags),
		),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Host region").
				Description("Please select the region of the host").
				OptionsFunc(func() []huh.Option[string] {
					regs, err := region.GetRegions(db)
					if err != nil {
						return []huh.Option[string]{
							huh.NewOption(
								"Error occurred while fetcing regions from database, please exit and retry",
								"-1",
							),
						}
					}
					if len(*regs) == 0 {
						return []huh.Option[string]{
							huh.NewOption(
								"No region present to use, please exit and create a region first",
								"-1",
							),
						}
					}
					opts := []huh.Option[string]{}
					for _, reg := range *regs {
						opts = append(opts, huh.NewOption(reg.Name, reg.Id.String()))
					}
					return opts
				}, nil).Value(&regionIDString).Validate(
				func(s string) error {
					if s == "-1" {
						return fmt.Errorf("no region present to use, please exit and create a region first")
					}
					return nil
				},
			),
		),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Identity File").
				Description("Please select the ssh identitity file to use for creating the connection with host").
				OptionsFunc(func() []huh.Option[string] {
					ids, err := identity.GetIdentity(db)
					if err != nil {
						return []huh.Option[string]{
							huh.NewOption(
								fmt.Sprintf("error occurred while fetching identities from database: %v", err), "-1",
							),
						}
					}

					if len(*ids) == 0 {
						return []huh.Option[string]{
							huh.NewOption(
								"No identites found in the database, please exit and insert identities first", "-1",
							),
						}
					}

					var opts []huh.Option[string]

					for _, id := range *ids {
						opts = append(opts,
							huh.NewOption(
								fmt.Sprintf("%s (%s)", id.Name, id.Path), id.Id.String(),
							),
						)
					}
					return opts

				}, nil).Value(&identityIDString).Validate(func(s string) error {
				if s == "-1" {
					return fmt.Errorf("no identites found in the database, please exit and insert identities first")
				}
				return nil
			}),
		),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Host Jumphost").
				Description("Please select the jumpost for the host, if any").
				OptionsFunc(func() []huh.Option[string] {
					hosts, err := GetShortHosts(db)
					if err != nil {
						return []huh.Option[string]{
							huh.NewOption(
								fmt.Sprintf("error occurred while fetchinig hosts: %v", err), "-1",
							),
						}
					}

					if len(*hosts) == 0 {
						return []huh.Option[string]{
							huh.NewOption("no host present to select", "0"),
						}
					}

					opts := []huh.Option[string]{
						huh.NewOption(
							"No jumphost", "0",
						),
					}
					for _, sh := range *hosts {
						opts = append(opts, huh.NewOption(
							sh.Name, sh.Id.String(),
						))
					}
					return opts

				}, nil).Value(&jumphostIDString).Validate(func(s string) error {
				if s == "-1" {
					return fmt.Errorf("error occurred while selecting jumphost, please exit and retry")
				}
				return nil
			}),
		),
	).WithTheme(huh.ThemeFunc(theme.XSH))

	if err := form.Run(); err != nil {
		return err
	}

	host.RegionID, _ = uuid.Parse(regionIDString)
	host.IdentityID, _ = uuid.Parse(identityIDString)

	if jumphostIDString != "0" {
		host.JumphostID.UUID = uuid.MustParse(jumphostIDString)
		host.JumphostID.Valid = true
	}
	host.Port, _ = strconv.Atoi(portString)

	host.Id, err = uuid.NewRandom()
	if err != nil {
		return err
	}

	return host.Store(db)
}
