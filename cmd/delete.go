package cmd

import (
	"database/sql"
	"xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [resource] [identifier]",
	Short: "Delete data from the database.",
	Long:  `Delete data from the database based on specified criteria.`,
}

var deleteHostCmd = &cobra.Command{
	Use:     "host [identifier]",
	Aliases: []string{"h"},
	Short:   "Delete host from the database.",
	Long: `Delete host from the database based on specified criteria.

Arguments:
  identifier: Any identifier for the resource selection. Please use * for selecting all
 `,
	RunE: func(cmd *cobra.Command, args []string) error {
		return genericDelete(args[0], host.Delete)
	},
}

var deleteRegionCmd = &cobra.Command{
	Use:     "region [identifier]",
	Aliases: []string{"r"},
	Short:   "Delete region from the database.",
	Long: `Delete region from the database based on specified criteria.

Arguments:
  identifier: Any identifier for the resource selection. Please use * for selecting all
 `,
	RunE: func(cmd *cobra.Command, args []string) error {
		return genericDelete(args[0], region.Delete)
	},
}

var deleteIdentityCmd = &cobra.Command{
	Use:     "identity [identifier]",
	Aliases: []string{"i"},
	Short:   "Delete identity from the database.",
	Long: `Delete identity from the database based on specified criteria.

Arguments:
  identifier: Any identifier for the resource selection. Please use * for selecting all
 `,
	RunE: func(cmd *cobra.Command, args []string) error {
		return genericDelete(args[0], identity.Delete)
	},
}

func genericDelete(identifer string, deletecFunc func(*sql.DB, string) error) error {
	dbConnection, err := db.GetDB()
	if err != nil {
		return err
	}

	defer dbConnection.Close()
	return deletecFunc(dbConnection, identifer)
}
