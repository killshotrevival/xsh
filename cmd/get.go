package cmd

import (
	"database/sql"
	"fmt"
	db "xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"

	"github.com/spf13/cobra"
)

var (
	outputFormat  string
	getIdentifier string
)

var getCmd = &cobra.Command{
	Use:   "get [resource]",
	Short: "Get data from the database.",
	Long: `Retrieve data from the database based on specified criteria.
	
Arguments:
  resource: Type of the resource. Possible values are (i)dentity / (h)ost / (r)egion
  identifier: Any identifier for the resource selection. Please use * for selecting all
	`,
}

var getHostCmd = &cobra.Command{
	Use:     "host",
	Aliases: []string{"h"},
	Short:   "Get host from the database.",
	Long: `Retrieve host from the database based on specified criteria.
	
Arguments:
  identifier: Any identifier for the resource selection. Please use * for selecting all
	`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return genericGetData(getIdentifier, outputFormat, host.Print)
	},
}

var getRegionCmd = &cobra.Command{
	Use:     "region",
	Aliases: []string{"r"},
	Short:   "Get region from the database.",
	Long: `Retrieve region from the database based on specified criteria.
	
Arguments:
  identifier: Any identifier for the resource selection. Please use * for selecting all
	`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return genericGetData(getIdentifier, outputFormat, region.Print)
	},
}

var getIdentityCmd = &cobra.Command{
	Use:     "identity",
	Aliases: []string{"i"},
	Short:   "Get identity from the database.",
	Long: `Retrieve identity from the database based on specified criteria.
	
Arguments:
  identifier: Any identifier for the resource selection. Please use * for selecting all
	`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return genericGetData(getIdentifier, outputFormat, identity.Print)
	},
}

func genericGetData(identifier, outputFormat string, getFunction func(*sql.DB, string, string) error) error {
	dbConnection, err := db.GetDB()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer dbConnection.Close()

	return getFunction(dbConnection, identifier, outputFormat)
}
