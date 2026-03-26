package region

import (
	"testing"
	"xsh/internal/db"
)

func TestPutRegion(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	if err := PutRegion(dbConnection, testResgionName); err != nil {
		t.Fatalf("error occurred while adding region to database: %v", err)
	}
}
