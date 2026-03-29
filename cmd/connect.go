// This file will hold code for cli to create ssh connection with server
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"

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

func sshConnect(_ *cobra.Command, args []string) error {
	sshString, err := buildSSHString(args[0])
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

func buildSSHString(identifier string) (string, error) {
	var (
		cHost, cjumpHost            *host.Host
		cIdentity, cJumhostIdentity *identity.Identity
	)
	dbConnection, err := db.GetDB()
	if err != nil {
		return "", err
	}

	cHost, err = host.GetHostByName(dbConnection, identifier)
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
