package host

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"xsh/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetHost(t *testing.T) {
	dbConnection := MockHostResources(t)
	defer dbConnection.Close()

	h, err := GetHostByName(dbConnection, testHostName)
	if err != nil {
		t.Fatalf("error occurred while fetching the host from database: %v", err)
	}

	assert.Equal(t, testHostName, h.Name)
}

func TestGetShortHosts(t *testing.T) {
	dbConnection := MockHostResources(t)
	defer dbConnection.Close()

	sh, err := GetShortHosts(dbConnection)
	if err != nil {
		t.Fatalf("error occurred while trying to fetch the hosts from database: %v", err)
	}

	assert.Equal(t, testJumpHostName, (*sh)[0].Name)

}

func TestPrintHost(t *testing.T) {
	dbConnection := MockHostResources(t)
	defer dbConnection.Close()

	path := utils.GetXSHTempDir(t)
	defer utils.RemoveTempDir(path, t)

	outputPath := filepath.Join(path, "test-host.json")

	err := Print(dbConnection, "*", "json", outputPath)
	assert.Nil(t, err)

	testData := []map[string]json.RawMessage{}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("error occurred while trying to read the host output file: %v", err)
	}

	if err := json.Unmarshal(data, &testData); err != nil {
		t.Fatalf("error occurred while trying to unmarshall the test data: %v", err)
	}

	if len(testData) == 0 {
		t.Fatalf("no data found in the host output file")
	}

}

func TestValidateFlags(t *testing.T) {
	err := validateExtraFlags("-J root@jumphost.com")
	assert.Error(t, err)

	err = validateExtraFlags("-4")
	assert.Nil(t, err)

}
