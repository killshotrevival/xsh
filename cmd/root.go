// Root file for the xsh command line
package cmd

import (
	"fmt"
	"xsh/internal/db"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/log"
)

var debug bool

// This variable is replaced in runtime while building the application
// for example: go build -ldflags "-X 'xsh/cmd.Version=1.2.3'"
var Version = "dev"

var rootCmd = &cobra.Command{
	Version: Version,
	Use:     "xsh",
	Short:   "Extended SSH",
	Long:    "A tool to extend the functionality of SSH with additional features and capabilities.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(true)
			log.SetReportTimestamp(true)
		}

		present, err := db.CheckDB()
		if err != nil {
			log.Debugf("error occurred while checking if DB file is present: %v", err)
			return nil
		}

		if present {
			log.Debug("DB found, checking for migrations to apply")
			if err := db.CheckAndApplyMigrations(); err != nil {
				return fmt.Errorf("error occurred while checking and applying database migration: %v", err)
			}
		} else {
			log.Debug("Datbase file not found, seems like the user is yet to init the environment.")
		}
		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Global flag available to all subcommands
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")

	rootCmd.AddCommand(putCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(connectCmd)
}
