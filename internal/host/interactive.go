package host

import (
	"database/sql"
	"fmt"
	"strconv"
	"xsh/internal/identity"
	"xsh/internal/region"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
	"github.com/google/uuid"
)

var (
	hostCreateOptions = map[int]string{
		1: "Clone from an existing host",
		2: "Create a new host from scratch",
	}
)

func xshTheme(isDark bool) *huh.Styles {
	t := huh.ThemeBase(isDark)
	lightDark := lipgloss.LightDark(isDark)

	var (
		normalFg = lightDark(lipgloss.Color("252"), lipgloss.Color("235"))
		indigo   = lightDark(lipgloss.Color("#5A56E0"), lipgloss.Color("#7571F9"))
		cream    = lightDark(lipgloss.Color("#FFFDF5"), lipgloss.Color("#FFFDF5"))
		fuchsia  = lipgloss.Color("#F780E2")
		green    = lightDark(lipgloss.Color("#02BA84"), lipgloss.Color("#02BF87"))
		red      = lightDark(lipgloss.Color("#FF4672"), lipgloss.Color("#ED567A"))
	)

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(indigo)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(indigo).MarginBottom(1)
	t.Focused.Directory = t.Focused.Directory.Foreground(indigo)
	t.Focused.Description = t.Focused.Description.Foreground(lightDark(lipgloss.Color(""), lipgloss.Color("243")))
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(fuchsia)
	t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(fuchsia)
	t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(fuchsia)
	t.Focused.Option = t.Focused.Option.Foreground(normalFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(fuchsia)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color("#02CF92"), lipgloss.Color("#02A877"))).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lightDark(lipgloss.Color(""), lipgloss.Color("243"))).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(normalFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(cream).Background(fuchsia)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(normalFg).Background(lightDark(lipgloss.Color("237"), lipgloss.Color("252")))

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(green)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lightDark(lipgloss.Color("248"), lipgloss.Color("238")))
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(fuchsia)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description
	return t
}

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
	).WithTheme(huh.ThemeFunc(xshTheme))

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
		cloneHostId string
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

				}, nil).Value(&cloneHostId).Validate(func(s string) error {
				if s == "-1" {
					return fmt.Errorf("error occurred while selecting host, please exit and retry")
				}
				return nil
			}),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	host, err := GetHostByID(db, cloneHostId)
	if err != nil {
		return err
	}

	host.Name = ""

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Host Name").
				Description("Please enter a unique name for your host that is easy to remember").
				Value(&host.Name).Validate(func(_ string) error {
				// TODO: Add support for validating that the given name does not exists
				return nil
			}),

			huh.NewInput().
				Title("Host Address").
				Description("Please enter the hostname / IP address of the host to connect").
				Value(&host.Address).Validate(func(_ string) error {
				// TODO: Add support for validating that the given address exists or not.
				return nil
			}),
		),
	)

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
				Value(&host.Name).Validate(func(_ string) error {
				// TODO: Add support for validating that the given name does not exists
				return nil
			}),

			huh.NewInput().
				Title("Host Address").
				Description("Please enter the hostname / IP address of the host to connect").
				Value(&host.Address).Validate(func(_ string) error {
				// TODO: Add support for validating that the given address exists or not.
				return nil
			}),

			huh.NewInput().
				Title("Host Port").
				Description("Please enter the port on which host accepts the connection").
				Value(&portString).Placeholder(portString),

			huh.NewInput().
				Title("Remote User").
				Description("Please enter the remote username").
				Value(&host.User).Placeholder(host.User),
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
	)

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

// func showError(app *tview.Application, msg string, previousPrimitive tview.Primitive) {
// 	modal := tview.NewModal().
// 		SetText(msg).
// 		AddButtons([]string{"OK"}).
// 		SetDoneFunc(func(i int, l string) {
// 			app.SetRoot(previousPrimitive, true) // restore list
// 		})

// 	app.SetRoot(modal, true)
// }

// func NewInteractivePut(db *sql.DB) error {
// 	app := tview.NewApplication()
// 	list := tview.NewList().
// 		AddItem("Clone", "Clone from an existring host", 'a', nil).
// 		AddItem("New", "Create a new host from scratch", 'b', nil).
// 		AddItem("Quit", "Press to exit", 'q', func() {
// 			app.Stop()
// 		})

// 	list = list.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
// 		defer app.Stop()
// 		if i == 0 {
// 			log.Debug("Cloning the host")
// 			showError(app, "Cloning it is", list)
// 			return
// 		}
// 		log.Debug("Creating a new host")
// 	})
// 	// list := tview.NewList().
// 	// 	AddItem("Option 1", "First choice", '1', func() {
// 	// 		fmt.Println("Selected Option 1")
// 	// 		app.Stop()
// 	// 	}).
// 	// 	AddItem("Option 2", "Second choice", '2', func() {
// 	// 		fmt.Println("Selected Option 2")
// 	// 		app.Stop()
// 	// 	}).
// 	// 	AddItem("Option 3", "Third choice", '3', func() {
// 	// 		fmt.Println("Selected Option 3")
// 	// 		app.Stop()
// 	// 	})

// 	// list.SetBorder(true).SetTitle("Choose an option")

// 	if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
// 		return err
// 	}

// 	return nil
// }
