package tool

import (
	"testing"
	"xsh/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()

	err := PutTool(dbConnection, "Testing", "testing")
	assert.Nil(t, err)

	tools, err := GetTools(dbConnection)
	assert.Nil(t, err)

	assert.Equal(t, "SSH", (*tools)[0].Name)
	assert.Equal(t, "Testing", (*tools)[1].Name)
}
