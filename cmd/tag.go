// Tagging resources are freezed until further development
package cmd

import (
	"database/sql"
	"xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"

	"github.com/spf13/cobra"
)

var remove bool

var tagCmd = &cobra.Command{
	Use:   "tag [command] [identifier] [tag]",
	Short: "Tag resources",
	Long: `Tag resources so that filtering can be done easily
	
Arguments:
  command: Type of the resource.
  identifier: Any identifier for the resource selection.
  tag: Tag value you want to place / remove on the resource
	`,
}

var tagHostCmd = &cobra.Command{
	Use:     "host [identifier] [tag]",
	Aliases: []string{"h"},
	Args:    cobra.ExactArgs(2),
	Short:   "Tag host",
	Long: `Tag host so that filtering can be done easily

Arguments:
  identifier: Any identifier for the resource selection.
  tag: Tag value you want to place / remove on the resource
	`,
	RunE: func(_ *cobra.Command, args []string) error {
		return genericTag(args[0], args[1], host.PutTagMapping, host.DeleteTagMapping)
	},
}

var tagRegionCmd = &cobra.Command{
	Use:     "region [identifier] [tag]",
	Aliases: []string{"r"},
	Short:   "Tag region",
	Args:    cobra.ExactArgs(2),
	Long: `Tag region so that filtering can be done easily

Arguments:
  identifier: Any identifier for the resource selection.
  tag: Tag value you want to place / remove on the resource
	`,
	RunE: func(_ *cobra.Command, args []string) error {
		return genericTag(args[0], args[1], region.PutTagMapping, region.DeleteTagMapping)
	},
}

var TagIDentityCmd = &cobra.Command{
	Use:     "identity [identifier] [tag]",
	Aliases: []string{"i"},
	Args:    cobra.ExactArgs(2),
	Short:   "Tag identity",
	Long: `Tag identity so that filtering can be done easily

Arguments:
  identifier: Any identifier for the resource selection.
  tag: Tag value you want to place / remove on the resource
	`,
	RunE: func(_ *cobra.Command, args []string) error {
		return genericTag(args[0], args[1], identity.PutTagMapping, identity.DeleteTagMapping)
	},
}

func genericTag(
	identifier, tagValue string,
	putTagMapping,
	deleteTagMapping func(*sql.DB, string, string) error,
) error {
	dbConnection, err := db.GetDB()
	if err != nil {
		return err
	}

	defer dbConnection.Close()
	if remove {
		return deleteTagMapping(dbConnection, identifier, tagValue)
	}
	return putTagMapping(dbConnection, identifier, tagValue)
}
