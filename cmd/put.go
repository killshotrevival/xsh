// This will hold code for cli to store data in database
package cmd

import (
	"fmt"
	"xsh/internal/identity"

	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:     "put",
	Short:   "Store data in the database.",
	Long:    "Store data in the database based on specified criteria.",
	Args:    cobra.ExactArgs(1),
	RunE:    putData,
	Example: "xsh put identity",
}

func putData(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "i":
		fallthrough
	case "identity":
		return identity.PutIdentity()
	default:
		return fmt.Errorf("invalid data type selected for inserting")
	}
}
