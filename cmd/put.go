// This will hold code for cli to store data in database
package cmd

import (
	"fmt"
	"os"
	db "xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"

	"github.com/spf13/cobra"
)

var putFile string
var interactivePut bool

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "Store data in the database.",
	Long:  "Store data in the database based on specified criteria.",
}

var putHostCmd = &cobra.Command{
	Use:     "host",
	Aliases: []string{"h"},
	Short:   "Store host in database",
	Long:    "Store host data in database either from a file or interactive prompting",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbConnection, err := db.GetDB()
		if err != nil {
			return fmt.Errorf("Error connecting to database: %w", err)
		}
		defer dbConnection.Close()
		if interactivePut {
			if putFile != putFileExample {
				return fmt.Errorf("interactive and file flags are mutually exclusive and please remove one flag and retry.")
			}
			return host.InteractivePut(dbConnection)
		}
		if _, err := os.Stat(putFile); err == os.ErrNotExist {
			return fmt.Errorf("%s file does not exists", putFile)
		}
		return host.PutHost(dbConnection, putFile)
	},
	Example: "xsh put host",
}

var putIdentityCmd = &cobra.Command{
	Use:     "identity",
	Aliases: []string{"i"},
	Short:   "Store SSH identity file in database",
	Long:    "Store SSH identity file in database which will be used for making connection with hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbConnection, err := db.GetDB()
		if err != nil {
			return fmt.Errorf("Error connecting to database: %w", err)
		}
		defer dbConnection.Close()

		name, path := args[0], args[1]
		return identity.PutIdentity(dbConnection, name, path)
	},
	Example: "xsh put identity [name] [path to private key file]",
}

var puRegionCmd = &cobra.Command{
	Use:     "region",
	Aliases: []string{"r"},
	Short:   "Store host regions in database",
	Long:    "Store host regions in database which will be used for better categorisation of hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbConnection, err := db.GetDB()
		if err != nil {
			return fmt.Errorf("Error connecting to database: %w", err)
		}
		defer dbConnection.Close()

		name := args[0]
		return region.PutRegion(dbConnection, name)
	},
	Example: "xsh region identity [name]",
}
