package host

import (
	"testing"
	"xsh/internal/db"
)

var (
	testHostName         = "test-host"
	testJumpHostName     = "jumphost-1"
	testHostJSONFilePath = "testdata/mock_host.json"
)

func TestPutHost(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutHost(dbConnection, testHostJSONFilePath); err != nil {
		t.Fatalf("error occurred while adding host to database: %v", err)
	}
}
