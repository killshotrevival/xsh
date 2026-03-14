package cmd

import (
	"fmt"
	"xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"

	"github.com/spf13/cobra"
)

var remove bool

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Tag resources",
	Long:  "Tag resources so that filtering can be done easily",
	Args:  cobra.ExactArgs(3),
	RunE:  tagData,
}

func tagData(cmd *cobra.Command, args []string) error {
	dataType, dataTypeValue, tagValue := args[0], args[1], args[2]

	dbConnection, err := db.GetDB()
	if err != nil {
		return err
	}
	switch dataType {
	case "i":
		fallthrough
	case "identity":
		if remove {
			return identity.DeleteTagMapping(dbConnection, dataTypeValue, tagValue)
		}
		return identity.PutTagMapping(dbConnection, dataTypeValue, tagValue)
	case "r":
		fallthrough
	case "region":
		if remove {
			return region.DeleteTagMapping(dbConnection, dataTypeValue, tagValue)
		}
		return region.PutTagMapping(dbConnection, dataTypeValue, tagValue)

	case "h":
		fallthrough
	case "host":
		if remove {
			return host.DeleteTagMapping(dbConnection, dataTypeValue, tagValue)
		}
		return host.PutTagMapping(dbConnection, dataTypeValue, tagValue)
	default:
		return fmt.Errorf("invalid data type selected for tagging")
	}
}
