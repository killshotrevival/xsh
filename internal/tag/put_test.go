package tag

import (
	"testing"
	"xsh/internal/db"
)

var (
	testTag = "test-tag=true"
)

func TestPutTag(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := Put(dbConnection, testTag); err != nil {
		t.Fatalf("error occurred while adding tag to database: %v", err)
	}
}
