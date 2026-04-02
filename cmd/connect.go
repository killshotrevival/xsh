// This file will hold code for cli to create ssh connection with server
package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/theme"

	"charm.land/huh/v2"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	printConnectionString bool
	verboseConnection     bool
)

var connectCmd = &cobra.Command{
	Use:   "connect [host name]",
	Short: "Connect SSH.",
	Long:  "Create an SSH connection with the specified server.",
	Args:  cobra.MaximumNArgs(2),
	RunE:  sshConnect,
}

func interactiveHostSelect(dbConnection *sql.DB) (string, error) {
	var h string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select host").
				Description("Please select the host you want to connect to").
				OptionsFunc(func() []huh.Option[string] {
					opts := []huh.Option[string]{}

					hosts, err := host.GetShortHosts(dbConnection)
					if err == nil {
						for _, ho := range *hosts {
							opts = append(opts, huh.NewOption(ho.Name, ho.Name))
						}
					} else {
						opts = append(opts, huh.NewOption("error occurred while fetching hosts from database, please try again with debug flag", "-1"))
					}

					return opts
				}, nil).
				Value(&h),
		),
	).WithTheme(huh.ThemeFunc(theme.XSH))

	if err := form.Run(); err != nil {
		return "", err
	}

	return h, nil
}

func sshConnect(_ *cobra.Command, args []string) error {
	var host string
	var err error

	dbConnection, err := db.GetDB()
	if err != nil {
		return err
	}

	if len(args) > 0 {
		host = args[0]
	} else {
		host, err = interactiveHostSelect(dbConnection)
		if err != nil {
			return err
		}
	}

	sshString, err := buildSSHString(host, dbConnection)
	if err != nil {
		return err
	}

	if verboseConnection {
		sshString += " -v"
	}

	if len(args) == 2 {
		sshString += fmt.Sprintf(" %s ", args[1])
	}

	if printConnectionString {
		log.Infof("[connect] Connecting to host with: %s", sshString)
		return nil
	}
	command := exec.Command("bash", "-c", sshString)

	// Attach terminal directly
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	err = command.Run()
	if err != nil {
		log.Errorf("There was some error observed while connecting over ssh: %v", err)
		log.Error("- If you are facing some issue, we are extremly sorry for that")
		log.Error("- If this execption occurred while exiting the ssh session, please ignore it. We are building a patch for this ")
		return err
	}
	log.Debug("[connect] SSH session completed successfully")
	return nil
}

func buildSSHString(identifier string, dbConnection *sql.DB) (string, error) {
	var (
		cHost, cjumpHost            *host.Host
		cIdentity, cJumhostIdentity *identity.Identity
	)

	cHost, err := host.GetHostByName(dbConnection, identifier)
	if err != nil {
		return "", err
	}

	cIdentity, err = identity.GetIdentityByID(dbConnection, cHost.IdentityID)
	if err != nil {
		return "", err
	}

	if cHost.JumphostID.Valid {
		cjumpHost, err = host.GetHostByID(dbConnection, cHost.JumphostID.UUID.String())
		if err != nil {
			return "", err
		}

		cJumhostIdentity, err = identity.GetIdentityByID(dbConnection, cjumpHost.IdentityID)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(`ssh -i %s -p %d -o ProxyCommand="ssh -i %s -W %s:%d %s@%s -p %d" %s@%s`,
			cIdentity.Path,
			cHost.Port,
			cJumhostIdentity.Path,
			cHost.Address,
			cHost.Port,
			cjumpHost.User,
			cjumpHost.Address,
			cjumpHost.Port,
			cHost.User,
			cHost.Address,
		), nil
	}

	return fmt.Sprintf("ssh -i %s -p %d %s@%s",
		cIdentity.Path,
		cHost.Port,
		cHost.User,
		cHost.Address,
	), nil
}
