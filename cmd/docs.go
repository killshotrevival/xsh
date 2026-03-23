package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var usageDoc bool

var docsCmd = &cobra.Command{
	Use:   "gendocs",
	Short: "Generate CLI documentation",
	RunE: func(cmd *cobra.Command, args []string) error {
		return doc.GenMarkdownTree(rootCmd, "./docs")
	},
}
