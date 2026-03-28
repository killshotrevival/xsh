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

var outputFile string

var exampleCmd = &cobra.Command{
	// TODO: Take output file name as input
	Use:   "example [resource]",
	Short: "Generate example format for resource",
	Long:  `Get JSON example of the datatype under a file. Can be used for JSON inserting data`,
}

var exampleHostCmd = &cobra.Command{
	Use:     "host",
	Aliases: []string{"h"},
	Short:   "Generate example format for host",
	Long:    `Get JSON example of the datatype under a file. Can be used for JSON inserting data`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return generateExampleFromComments(host.Host{}, []string{"id", "region_name", "jumphost_name", "identitiy_file_name", "tags"}, outputFile)
	},
}

var exampleIdentityCmd = &cobra.Command{
	Use:     "identity",
	Aliases: []string{"i"},
	Short:   "Generate example format for identity",
	Long:    `Get JSON example of the datatype under a file. Can be used for JSON inserting data`,
	RunE: func(_ *cobra.Command, _ []string) error {
		return generateExampleFromComments(identity.Identity{}, []string{"id", "tags"}, outputFile)
	},
}

func generateExampleFromComments(v any, ignoreKeys []string, output string) error {
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

	return os.WriteFile(output, []byte(result), 0644)

}
