package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:   "gendocs",
	Short: "Generate CLI documentation",
	RunE: func(_ *cobra.Command, _ []string) error {
		return doc.GenMarkdownTree(rootCmd, "./docs")
	},
}
