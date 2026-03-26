package region

import (
	"strings"
	"testing"
	"xsh/internal/db"
)

func TestDeleteRegion(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutRegion(dbConnection, testResgionName); err != nil {
		t.Fatalf("error occurred while adding region to database: %v", err)
	}

	if err := Delete(dbConnection, testResgionName); err != nil {
		t.Fatalf("error occurre while deleting the region: %v", err)
	}

	reg, err := GetRegionByName(dbConnection, testResgionName)
	if err != nil {
		if strings.Contains(err.Error(), noRegionErr) {
			return
		}
		t.Fatalf("error occurred while trying to fetch the region with name: %v", err)
	}
	t.Fatalf("test region found in database, even after successful deletion: %s", reg)
}
