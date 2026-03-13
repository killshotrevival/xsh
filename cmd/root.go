// Root file for the xsh command line
package cmd

import (
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
			log.SetReportCaller(true)
			log.SetReportTimestamp(true)
		}
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
