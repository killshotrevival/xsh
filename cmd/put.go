// This will hold code for cli to store data in database
package cmd

import (
	"fmt"
	db "xsh/internal/db"
	"xsh/internal/identity"

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

	dataType, name, path := args[0], args[1], args[2]
	switch dataType {
	case "i":
		fallthrough
	case "identity":
		return identity.PutIdentity(dbConnection, name, path)
	default:
		return fmt.Errorf("invalid data type selected for inserting")
	}
}
