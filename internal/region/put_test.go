package region

import (
	"testing"
	"xsh/internal/db"
)

var (
	testResgionName = "testing-region"
)

func TestPutRegion(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutRegion(dbConnection, testResgionName); err != nil {
		t.Fatalf("error occurred while adding region to database: %v", err)
	}
}
