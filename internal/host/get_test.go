package host

import (
	"testing"
	"xsh/internal/db"

	"github.com/stretchr/testify/assert"
)

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

	assert.Equal(t, h.Name, testHostName)
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

	assert.Equal(t, (*sh)[0].Name, testHostName)

}
