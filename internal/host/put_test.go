package host

import (
	"database/sql"
	"testing"
	"xsh/internal/db"
	"xsh/internal/identity"
	"xsh/internal/region"
	"xsh/internal/tool"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	testHostName         = "test-host"
	testJumpHostName     = "jumphost-1"
	testHostJSONFilePath = "testdata/mock_host.json"
)

func mockRegionData() *region.Region {
	return &region.Region{
		Id:   uuid.MustParse("dbd31312-f575-4d2f-803d-57b506b14d0c"),
		Name: "us-east-1",
	}
}

func mockIdentityData() *identity.Identity {
	return &identity.Identity{
		Id:   uuid.MustParse("7f1620cb-b627-4f32-a7eb-e05d19d6065c"),
		Name: "test-identity",
		Path: "/Users/test/.ssh/test_key",
	}
}

func MockHostResources(t *testing.T) *sql.DB {
	dbConnection := db.GetTestDB(t)

	r := mockRegionData()
	if err := r.Store(dbConnection); err != nil {
		t.Fatalf("error occurred while storing mock region data: %v", err)
	}

	i := mockIdentityData()
	if err := i.Store(dbConnection); err != nil {
		t.Fatalf("error occurred while storing mock identity data: %v", err)
	}
	if err := PutHost(dbConnection, testHostJSONFilePath); err != nil {
		t.Fatalf("error occurred while adding host to database: %v", err)
	}

	return dbConnection
}

func TestPutHost(t *testing.T) {
	db := MockHostResources(t)
	defer db.Close()

	h, err := GetHostByName(db, "jumphost-1")
	assert.Nil(t, err)

	to, err := tool.GetToolByID(db, h.ToolID)
	assert.Nil(t, err)

	assert.Equal(t, "SSH", to.Name)

}
