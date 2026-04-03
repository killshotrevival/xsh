// Root file for the xsh command line
package cmd

import (
	"fmt"
	"xsh/internal/db"
	"xsh/internal/theme"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/log"
)

var (
	debug bool

	// This variable is replaced in runtime while building the application
	// for example: go build -ldflags "-X 'xsh/cmd.Version=1.2.3'"
	Version = "dev"

	putFileExample = "/path/to/file"
)

var rootCmd = &cobra.Command{
	Version: Version,
	Use:     "xsh",
	Short:   "Extended SSH",
	Long:    "A tool to extend the functionality of SSH with additional features and capabilities.",
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		if debug {
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(true)
			log.SetReportTimestamp(true)
		}

		present, err := db.CheckDB()
		if err != nil {
			log.Debugf("[root] failed to verify database file existence: %v", err)
			return nil
		}

		if present {
			log.Debug("[root] database file found, checking for pending migrations")
			if err := db.CheckAndApplyMigrations(); err != nil {
				return fmt.Errorf("error occurred while checking and applying database migration: %v", err)
			}
		} else {
			log.Warn("[root] database file not found; run 'xsh init' to set up the environment")
		}

		// Apply tview theme
		theme.ApplyTviewTheme()

		return nil
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// Global flag available to all subcommands
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")
	rootCmd.AddCommand(docsCmd, putCmd, getCmd, initCmd, connectCmd, deleteCmd, exampleCmd, importCmd, editCmd)

	// Connect Command
	connectCmd.PersistentFlags().BoolVarP(&printConnectionString, "print", "p", false, "Just print connection string, instead of connecting with the host")
	connectCmd.PersistentFlags().BoolVarP(&verboseConnection, "verbose", "v", false, "Create ssh connection in verbose mode")

	// Import Command
	importCmd.PersistentFlags().BoolVarP(&printSourceFiles, "print", "p", false, "Print all the predefined files, XSH will be reading and importing from")

	// Put Command
	putCmd.AddCommand(putHostCmd, puRegionCmd, putIdentityCmd)
	putHostCmd.PersistentFlags().BoolVarP(&interactivePut, "interactive", "i", true, "Insert host in interactive mode")
	putHostCmd.PersistentFlags().StringVarP(&putFile, "file", "f", putFileExample, "Path of the host file")

	// Get Command
	getCmd.AddCommand(getHostCmd, getRegionCmd, getIdentityCmd)
	getCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "Format of the output")
	getCmd.PersistentFlags().StringVarP(&getIdentifier, "identifier", "i", "*", "Identifier for filtering the data")
	getCmd.PersistentFlags().StringVarP(&getOutputFile, "output-file", "f", "resource.json", "Output file name, only used when output format is `json`")

	// Delete Command
	deleteCmd.PersistentFlags().BoolVarP(&interactiveDelete, "interactive", "i", true, "Insert host in interactive mode")

	// Example Command
	exampleCmd.AddCommand(exampleHostCmd, exampleIdentityCmd)
	exampleCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "example.json", "Output file name")

	// Tag Command
	tagCmd.AddCommand(tagHostCmd, TagIDentityCmd, tagRegionCmd)
	tagCmd.PersistentFlags().BoolVarP(&remove, "remove", "r", false, "Remove Tag mapping")
}
