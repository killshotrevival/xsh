package identity

import (
	"path/filepath"
	"testing"
	"xsh/internal/db"

	"github.com/stretchr/testify/assert"
)

var (
	testIdentityName    = "test-id"
	testIdentityPath, _ = filepath.Abs("testdata/mock_identity_file")
)

func TestPutIdentity(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()
	if err := PutIdentity(dbConnection, testIdentityName, testIdentityPath); err != nil {
		t.Fatalf("error occurred while adding identity to database: %v", err)
	}
}

func TestPutIdentityPath(t *testing.T) {
	err := PutIdentity(nil, testIdentityName, "~/.ssh/testing")
	assert.Equal(t, errRelativeFilePath, err)

	err = PutIdentity(nil, testIdentityName, "/User/twelcon/Desktop/testing")
	assert.Contains(t, err.Error(), "no such file or directory")
}
