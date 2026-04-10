package tool

import (
	"testing"
	"xsh/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestPutTool(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()

	err := PutTool(dbConnection, "Testing", "testing")
	assert.Nil(t, err)
}
