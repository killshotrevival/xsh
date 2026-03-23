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

/*
TODO:
1. Create SSH connection with the server using the ssh cli of the system
2. Talk to the database and extract all the relivant information for the server we are connecting to and use it to enhance the ssh connection
*/

var connectCmd = &cobra.Command{
	Use:   "connect [host name]",
	Short: "Connect SSH.",
	Long:  "Create an SSH connection with the specified server.",
	Args:  cobra.ExactArgs(1),
	RunE:  sshConnect,
}

func sshConnect(_ *cobra.Command, args []string) error {
	sshString, err := buildSSHString(args[0])
	if err != nil {
		return err
	}

	fmt.Println(sshString)
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
	log.Info("Completed SSH session gracefully")
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
		return fmt.Sprintf(`ssh -i %s -o ProxyCommand="ssh -i %s -W %s:%d %s@%s" %s@%s`,
			cIdentity.Path,
			cJumhostIdentity.Path,
			cHost.Address,
			cHost.Port,
			cjumpHost.User,
			cjumpHost.Address,
			cHost.User,
			cHost.Address,
		), nil
	}

	return fmt.Sprintf("ssh -i %s %s@%s:%d",
		cIdentity.Path,
		cHost.User,
		cHost.Address,
		cHost.Port,
	), nil
}
