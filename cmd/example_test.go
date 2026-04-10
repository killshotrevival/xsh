package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"xsh/internal/host"
	"xsh/internal/utils"
)

func TestGenerateExampleFromComments(t *testing.T) {

	path := utils.GetXSHTempDir(t)

	defer utils.RemoveTempDir(path, t)

	output := filepath.Join(path, "example.json")

	if err := generateExampleFromComments(reflect.TypeFor[host.Host](), []string{"id", "region_name", "jumphost_name", "identitiy_file_name", "tags"}, output); err != nil {
		t.Fatalf("error occurred while creating example docs for host: %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("error occurred while reading the example file: %v", err)
	}

	var exampleData map[string]json.RawMessage

	if err := json.Unmarshal(data, &exampleData); err != nil {
		t.Fatalf("error occurred while unmarshalling the example data: %v", err)
	}

	if _, present := exampleData["name"]; !present {
		t.Fatalf("name key not present in the example data")
	}

	if _, present := exampleData["id"]; present {
		t.Fatalf("id present in the example data, when it should not be")
	}

}
