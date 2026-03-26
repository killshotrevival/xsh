// This file will hold code for cli to initialize the XSH environment
package cmd

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	config "xsh/internal/config"
	db "xsh/internal/db"
	"xsh/internal/identity"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the XSH environment.",
	Long:  "Set up the XSH environment by creating necessary directories and database tables.",
	RunE:  initXSH,
}

func initXSH(_ *cobra.Command, _ []string) error {

	log.Info("Iitialising xsh configuration directory")
	if err := config.InitConfigDir(); err != nil {
		log.Debugf("[init] failed to initialize configuration directory: %v", err)
		return err
	}

	log.Info("Initialising database")
	dbExists, err := db.CheckDB()
	if err != nil {
		return fmt.Errorf("error checking database: %w", err)
	}

	if !dbExists {
		err := db.InitDB()
		if err != nil {
			return fmt.Errorf("error initializing database: %w", err)
		}
		log.Info("Database initialized successfully.")
	}
	log.Info("Database already exists. Skipping initialization.")

	log.Info("Pre populating identity table")
	dbConnection, err := db.GetDB()
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	defer dbConnection.Close()

	err = identity.InitIdentityStore(dbConnection)
	if err != nil {
		return fmt.Errorf("error initializing identity store: %w", err)
	}
	log.Info("Identity table populated successfully")
	return nil
}
