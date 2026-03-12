// This will hold code for cli to store data in database
package cmd

import (
	"fmt"
	db "xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"

	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:     "put",
	Short:   "Store data in the database.",
	Long:    "Store data in the database based on specified criteria.",
	Args:    cobra.ExactArgs(3),
	RunE:    putData,
	Example: "xsh put identity",
}

func putData(cmd *cobra.Command, args []string) error {
	dbConnection, err := db.GetDB()
	if err != nil {
		return fmt.Errorf("Error connecting to database: %w", err)
	}
	defer dbConnection.Close()

	dataType := args[0]
	switch dataType {
	case "i":
		fallthrough
	case "identity":
		name, path := args[1], args[2]
		return identity.PutIdentity(dbConnection, name, path)
	case "r":
		fallthrough
	case "region":
		name := args[1]
		return region.PutRegion(dbConnection, name)
	case "h":
		fallthrough
	case "host":
		filepath := args[1]
		return host.PutHost(dbConnection, filepath)
	default:
		return fmt.Errorf("invalid data type selected for inserting")
	}
}
