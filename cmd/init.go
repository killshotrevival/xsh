// This file will hold code for cli to initialize the XSH environment
package cmd

/*
TODO:
1. Init the XSH root directory in the user's home directory (e.g., ~/.xsh)
2. Check if SQLite database exists, if not create it
3. Check if tables exists in the database, if not create them
4. Read ~/.ssh/config, ~/.bashrc, ~/.zshrc and store relevant information in the database
*/

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

func initXSH(cmd *cobra.Command, args []string) error {

	log.Info("Iitialising xsh configuration directory")
	if err := config.InitConfigDir(); err != nil {
		log.Debugf("Error occurred while initialising config dir: %v", err)
		return err
	}

	log.Info("Initialising database")
	dbExists, err := db.CheckDB()
	if err != nil {
		return fmt.Errorf("Error checking database: %w", err)
	}

	if !dbExists {
		err := db.InitDB()
		if err != nil {
			return fmt.Errorf("Error initializing database: %w", err)
		}
		log.Info("Database initialized successfully.")
	}
	log.Info("Database already exists. Skipping initialization.")

	log.Info("Pre populating identity table")
	dbConnection, err := db.GetDB()
	if err != nil {
		return fmt.Errorf("Error connecting to database: %w", err)
	}
	defer dbConnection.Close()

	err = identity.InitIdentityStore(dbConnection)
	if err != nil {
		return fmt.Errorf("Error initializing identity store: %w", err)
	}
	log.Info("Identity table populated successfully")
	return nil
}
