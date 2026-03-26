package identity

import (
	"testing"
	"xsh/internal/db"
)

var (
	testIdentityName = "test-id"
	testIdentityPath = "/testing"
)

func TestPutIdentity(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutIdentity(dbConnection, testIdentityName, testIdentityPath); err != nil {
		t.Fatalf("error occurred while adding identity to database: %v", err)
	}
}
