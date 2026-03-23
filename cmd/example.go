package cmd

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"xsh/internal/host"
	"xsh/internal/identity"

	"github.com/spf13/cobra"
)

var exampleCmd = &cobra.Command{
	Use:   "example [resource]",
	Short: "Generate example format for resource",
	Long:  `Get JSON example of the datatype under example.json file. Can be used for JSON inserting data`,
}

var exampleHostCmd = &cobra.Command{
	Use:     "host",
	Aliases: []string{"h"},
	Short:   "Generate example format for host",
	Long:    `Get JSON example of the datatype under example.json file. Can be used for JSON inserting data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateExampleFromComments(host.Host{}, []string{"id", "region_name", "jumphost_name", "identitiy_file_name", "tags"})
	},
}

var exampleIdentityCmd = &cobra.Command{
	Use:     "identity",
	Aliases: []string{"i"},
	Short:   "Generate example format for identity",
	Long:    `Get JSON example of the datatype under example.json file. Can be used for JSON inserting data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateExampleFromComments(identity.Identity{}, []string{"id", "tags"})
	},
}

func generateExampleFromComments(v any, ignoreKeys []string) error {
	val := reflect.TypeOf(v)

	result := "{\n"
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		jsonKey := strings.Split(field.Tag.Get("json"), ",")[0]

		if slices.Contains(ignoreKeys, jsonKey) {
			continue
		}
		comment := field.Tag.Get("comment")

		result += fmt.Sprintf("  \"%s\": \"%s\"", jsonKey, comment)

		if i < val.NumField()-1 {
			result += ","
		}
		result += "\n"
	}
	// Removing unnecessary "," at the end of the last line
	result = result[0 : len(result)-2]
	result += "\n}\n"

	return os.WriteFile("example.json", []byte(result), 0644)

}
