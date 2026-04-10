package cmd

import (
	"database/sql"
	"fmt"
	"strings"
	"xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"
	"xsh/internal/tool"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var interactiveDelete bool

var resourceDeleteMapping = map[string]func(*sql.DB, string) error{
	"host":     host.Delete,
	"region":   region.Delete,
	"identity": identity.Delete,
	"tool":     tool.Delete,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete data from the database.",
	Long:  `Delete data from the database based on specified criteria.`,
	Args:  cobra.MaximumNArgs(2),
	RunE:  deleteData,
}

func deleteData(_ *cobra.Command, args []string) error {
	dbConnection, err := db.GetDB()
	if err != nil {
		return err
	}
	defer dbConnection.Close()

	if len(args) == 2 {
		resource, identifier := args[0], args[1]
		for key, deleteFunc := range resourceDeleteMapping {
			if strings.Contains(key, strings.ToLower(resource)) {
				return deleteFunc(dbConnection, identifier)
			}
		}
		return fmt.Errorf("invalid resource type provided for deletion")
	}

	return Interactive(dbConnection)

}

func Interactive(dbConnection *sql.DB) error {

	resource, idLists, err := selectResource(dbConnection)
	if err != nil {
		return nil
	}

	for _, id := range idLists {
		log.Debugf("[delete] removing %s resource %q", resource, id)
		if err := resourceDeleteMapping[resource](dbConnection, id); err != nil {
			log.Warnf("[delete] failed to delete %s resource %q: %v", resource, id, err)
		}
	}

	return nil
}
