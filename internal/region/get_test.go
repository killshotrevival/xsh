package region

import (
	"testing"
	"xsh/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestGetRegionByName(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutRegion(dbConnection, testResgionName); err != nil {
		t.Fatalf("error occurred while adding region to database: %v", err)
	}

	reg, err := GetRegionByName(dbConnection, testResgionName)
	if err != nil {
		t.Fatalf("unable to fetch region by name: %v", err)
	}

	assert.Equal(t, reg.Name, testResgionName)
}

func TestGetRegions(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutRegion(dbConnection, testResgionName); err != nil {
		t.Fatalf("error occurred while adding region to database: %v", err)
	}

	regions, err := GetRegions(dbConnection)
	if err != nil {
		t.Fatalf("error occurred while trying to fetch regions from database: %v", err)
	}

	for _, reg := range *regions {
		if testResgionName == reg.Name {
			return
		}
	}

	t.Fatalf("%s region not found in database even though the put region succeeded", testResgionName)

}
