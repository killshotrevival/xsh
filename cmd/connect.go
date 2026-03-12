// This file will hold code for cli to create ssh connection with server
package cmd

import "github.com/spf13/cobra"

/*
TODO:
1. Create SSH connection with the server using the ssh cli of the system
2. Talk to the database and extract all the relivant information for the server we are connecting to and use it to enhance the ssh connection
*/

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect SSH.",
	Long:  "Create an SSH connection with the specified server.",
	RunE:  sshConnect,
}

func sshConnect(cmd *cobra.Command, args []string) error {
	return nil
}
