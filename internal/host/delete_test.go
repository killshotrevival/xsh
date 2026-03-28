package host

import (
	"strings"
	"testing"
	"xsh/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutHost(dbConnection, testHostJSONFilePath); err != nil {
		t.Fatalf("error occurred while adding host to database: %v", err)
	}

	if err := Delete(dbConnection, testHostName); err != nil {
		t.Fatalf("error occurred while trying to delete a host: %v", err)
	}

	if _, err := GetHostByName(dbConnection, testHostName); err != nil {
		if strings.Contains(err.Error(), noHostFoundError) {
			return
		}
		t.Fatalf("error occurred while trying to fetch host from database: %v", err)
	}

	t.Fatal("host present in database event though delete executed successfully")
}

func TestDeleteJumphost(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutHost(dbConnection, testHostJSONFilePath); err != nil {
		t.Fatalf("error occurred while adding host to database: %v", err)
	}

	err := Delete(dbConnection, testJumpHostName)

	assert.Equal(t, jumphostDeleteError, err)
}
