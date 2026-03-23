package cmd

import (
	"fmt"
	db "xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"
	"xsh/internal/tag"

	"github.com/spf13/cobra"
)

var (
	outputFormat string
)

var getCmd = &cobra.Command{
	Use:   "get [resource] [identifier]",
	Short: "Get data from the database.",
	Long: `Retrieve data from the database based on specified criteria.
	
Arguments:
  resource: Type of the resource. Possible values are (i)dentity / (h)ost / (r)egion
  identifier: Any identifier for the resource selection. Please use * for selecting all
	`,
	Args: cobra.ExactArgs(2),
	RunE: getData,
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
		return identity.Print(dbConnection, args[1], outputFormat)
	case "r":
		fallthrough
	case "region":
		_, err := region.PrintRegions(dbConnection, args[1], outputFormat)
		return err
	case "h":
		fallthrough
	case "hosts":
		return host.Print(dbConnection, args[1], outputFormat)
	case "t":
		fallthrough
	case "tag":
		return tag.Print(dbConnection, args[1], outputFormat)
	default:
		return fmt.Errorf("invalid data type selected for fetcing")
	}
}
