package tool

import (
	"testing"
	"xsh/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()

	to, err := NewTool("testing", "testing connnectio string")
	assert.Nil(t, err)

	err = to.Store(dbConnection)
	assert.Nil(t, err)

	to, err = GetToolByName(dbConnection, "testing")
	assert.Nil(t, err)

	err = Delete(dbConnection, to.ID.String())
	assert.Nil(t, err)

}
