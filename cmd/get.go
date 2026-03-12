// This will hold code for cli to get data from database
package cmd

import (
	"fmt"
	db "xsh/internal/db"
	"xsh/internal/identity"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get data from the database.",
	Long:  "Retrieve data from the database based on specified criteria.",
	Args:  cobra.ExactArgs(2),
	RunE:  getData,
}

func getData(cmd *cobra.Command, args []string) error {
	dbConnection, err := db.GetDB()
	if err != nil {
		return fmt.Errorf("Error connecting to database: %w", err)
	}
	defer dbConnection.Close()

	switch args[0] {
	case "i":
		fallthrough
	case "identity":
		return identity.PrintIdentities(dbConnection, args[1])
	default:
		return fmt.Errorf("invalid data type selected for fetcing")
	}
}
