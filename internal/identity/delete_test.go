package identity

import (
	"strings"
	"testing"
	"xsh/internal/db"
)

func TestDelete(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutIdentity(dbConnection, testIdentityName, testIdentityPath); err != nil {
		t.Fatalf("error occurred while adding region to database: %v", err)
	}

	err := Delete(dbConnection, testIdentityName)
	if err != nil {
		t.Fatalf("error occurred while deleting the identity: %v", err)
	}

	_, err = GetIdentityByName(dbConnection, testIdentityName)
	if err != nil {
		if strings.Contains(err.Error(), noIdentityFoundErr) {
			return
		}

		t.Fatalf("error occurred while trying to fetch the identity: %v", err)
	}

	t.Fatal("Identity still present in database even after successful delete operation")
}
