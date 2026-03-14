package cmd

import (
	"fmt"
	"xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete data from the database.",
	Long:  "Delete data from the database based on specified criteria.",
	Args:  cobra.ExactArgs(2),
	RunE:  deleteData,
}

func deleteData(cmd *cobra.Command, args []string) error {
	dataType, identifer := args[0], args[1]

	dbConnection, err := db.GetDB()
	if err != nil {
		return err
	}

	switch dataType {
	case "i":
		fallthrough
	case "identity":
		return identity.Delete(dbConnection, identifer)
	case "h":
		fallthrough
	case "host":
		return host.Delete(dbConnection, identifer)
	case "r":
		fallthrough
	case "regions":
		return region.Delete(dbConnection, identifer)
	default:
		return fmt.Errorf("invalid datatype selected for deletion")
	}
}
