package tag

import (
	"testing"
	"xsh/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestGetTag(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := Put(dbConnection, testTag); err != nil {
		t.Fatalf("error occurred while adding tag to database: %v", err)
	}

	tag, err := GetTag(dbConnection, testTag)
	if err != nil {
		t.Fatalf("error occurred while fetching tag from database: %v", err)
	}

	assert.Equal(t, testTag, tag.Tag)
}

func TestGetTagWithCreate(t *testing.T) {
	tempTag := "testing-Tag-Create"
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := Put(dbConnection, testTag); err != nil {
		t.Fatalf("error occurred while adding tag to database: %v", err)
	}

	_, err := GetTagWithCreate(dbConnection, tempTag)
	if err != nil {
		t.Fatalf("error occurred while fetching tag with create from database: %v", err)
	}

	tag, err := GetTag(dbConnection, tempTag)
	if err != nil {
		t.Fatalf("error occurred while fetching tag from database: %v", err)
	}

	assert.Equal(t, tempTag, tag.Tag)
}
