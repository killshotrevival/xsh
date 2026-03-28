package identity

import (
	"testing"
	"xsh/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestGetIdentity(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutIdentity(dbConnection, testIdentityName, testIdentityPath); err != nil {
		t.Fatalf("error occurred while adding region to database: %v", err)
	}

	ids, err := GetIdentity(dbConnection)
	if err != nil {
		t.Fatalf("error occurred while fetching identities from database: %v", err)
	}

	for _, id := range *ids {
		if id.Name == testIdentityName && id.Path == testIdentityPath {
			return
		}
	}
	t.Fatal("Identity not found in database even after successful put")
}

func TestGetIdentityByName(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutIdentity(dbConnection, testIdentityName, testIdentityPath); err != nil {
		t.Fatalf("error occurred while adding region to database: %v", err)
	}

	id, err := GetIdentityByName(dbConnection, testIdentityName)
	if err != nil {
		t.Fatalf("error occurred while trying to fetch identity with name: %v", err)
	}
	assert.Equal(t, testIdentityName, id.Name)
}
