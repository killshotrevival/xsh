package host

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"xsh/internal/db"
	"xsh/internal/identity"
	"xsh/internal/region"
	"xsh/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func mockRegionData() *region.Region {
	return &region.Region{
		Id:   uuid.MustParse("dbd31312-f575-4d2f-803d-57b506b14d0c"),
		Name: "us-east-1",
	}
}

func mockIdentityData() *identity.Identity {
	return &identity.Identity{
		Id:   uuid.MustParse("7f1620cb-b627-4f32-a7eb-e05d19d6065c"),
		Name: "test-identity",
		Path: "/Users/test/.ssh/test_key",
	}
}

func TestGetHost(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutHost(dbConnection, testHostJSONFilePath); err != nil {
		t.Fatalf("error occurred while adding host to database: %v", err)
	}

	h, err := GetHostByName(dbConnection, testHostName)
	if err != nil {
		t.Fatalf("error occurred while fetching the host from database: %v", err)
	}

	assert.Equal(t, testHostName, h.Name)
}

func TestGetShortHosts(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutHost(dbConnection, testHostJSONFilePath); err != nil {
		t.Fatalf("error occurred while adding host to database: %v", err)
	}

	sh, err := GetShortHosts(dbConnection)
	if err != nil {
		t.Fatalf("error occurred while trying to fetch the hosts from database: %v", err)
	}

	assert.Equal(t, testHostName, (*sh)[0].Name)

}

func TestPrintHost(t *testing.T) {

	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutHost(dbConnection, testHostJSONFilePath); err != nil {
		t.Fatalf("error occurred while adding host to database: %v", err)
	}

	r := mockRegionData()

	if err := r.Store(dbConnection); err != nil {
		t.Fatalf("error occurred while storing mock region data: %v", err)
	}

	i := mockIdentityData()

	if err := i.Store(dbConnection); err != nil {
		t.Fatalf("error occurred while storing mock identity data: %v", err)
	}

	path := utils.GetXSHTempDir(t)

	defer utils.RemoveTempDir(path, t)

	outputPath := filepath.Join(path, "test-host.json")

	if err := Print(dbConnection, "*", "json", outputPath); err != nil {
		t.Fatalf("error occurred while trying to print the hosts: %v", err)
	}

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
