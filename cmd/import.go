package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"xsh/internal/db"
	import_xsh "xsh/internal/import"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	printSourceFiles bool
	sourceFiles      = []string{
		".bash_history",
		".zsh_history",
		".zshrc",
		".bashrc",
		// ".ssh/config", TODO
	}
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import SSH configurations from predefined files",
	RunE:  importFromSourceFiles,
}

func importFromSourceFiles(_ *cobra.Command, _ []string) error {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error occurred while trying to get the home directory: %v", err)
	}

	dbConnection, err := db.GetDB()
	if err != nil {
		return err
	}

	for _, file := range sourceFiles {
		path := filepath.Join(homeDir, file)
		if printSourceFiles {
			log.Infof("Will be importing from: %s", path)
			continue
		}

		if err := import_xsh.Import(path, dbConnection); err != nil {
			log.Warnf("error occurred while importing %s: %v", path, err)
		}
	}

	return nil
}
