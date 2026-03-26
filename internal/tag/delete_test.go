package tag

import (
	"strings"
	"testing"
	"xsh/internal/db"
)

func TestDelete(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := Put(dbConnection, testTag); err != nil {
		t.Fatalf("error occurred while adding tag to database: %v", err)
	}

	if err := Delete(dbConnection, testTag); err != nil {
		t.Fatalf("error occurred while deleting the tag from database: %v", err)
	}

	_, err := GetTag(dbConnection, testTag)
	if err != nil {
		if strings.Contains(err.Error(), noTagFoundError) {
			return
		}
		t.Fatalf("error occurred while trying to fetch tag from database: %v", err)
	}
	t.Fatal("tag still present in the database, even after successful deletion")
}
