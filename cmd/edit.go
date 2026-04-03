package cmd

import (
	"database/sql"
	"fmt"
	"xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"
	"xsh/internal/region"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit resources",
	RunE:  editResources,
}

var regionEditFucntionMapping = map[string]func(*sql.DB, string) error{
	"host":     host.Edit,
	"region":   region.Edit,
	"identity": identity.Edit,
}

func editResources(_ *cobra.Command, _ []string) error {
	dbConnection, err := db.GetDB()
	if err != nil {
		return err
	}
	defer dbConnection.Close()
	resource, idLists, err := selectResource(dbConnection)
	if err != nil {
		return err
	}
	if len(idLists) > 1 {
		return fmt.Errorf("please select only one resource for editing")
	}

	return regionEditFucntionMapping[resource](dbConnection, idLists[0])
}
