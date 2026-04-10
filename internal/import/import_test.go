package import_xsh

import (
	"testing"
	"xsh/internal/db"
	"xsh/internal/host"
	"xsh/internal/identity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func MockIdentityData() *identity.Identity {
	return &identity.Identity{
		Id:   uuid.MustParse("7f1620cb-b627-4f32-a7eb-e05d19d6065c"),
		Name: "test-identity",
		Path: "/Users/test/.ssh/test_key",
	}
}

func TestCommandToHost(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()

	id := MockIdentityData()

	err := id.Store(dbConnection)
	assert.Nil(t, err)

	command := `ssh -4 -i /Users/test/.ssh/test_key root@test.com -p 202`

	h, err := commandToHost(command, dbConnection)
	assert.Nil(t, err)

	assert.Equal(t, "test.com", h.Name)
	assert.Equal(t, uuid.MustParse("7f1620cb-b627-4f32-a7eb-e05d19d6065c"), h.IdentityID)
	assert.Equal(t, 202, h.Port)

}

func TestCommandToHostWithProxyCommand(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()

	id := MockIdentityData()

	err := id.Store(dbConnection)
	assert.Nil(t, err)

	command := `ssh -4 -i /Users/test/.ssh/test_key -o ProxyCommand="ssh -i /Users/test/.ssh/test_key -W test.com:202 test@jumpbox.com -p 203" root@test.com -p 2002 -v`

	h, err := commandToHost(command, dbConnection)
	assert.Nil(t, err)

	assert.Equal(t, "test.com", h.Name)
	assert.Equal(t, 2002, h.Port)

	j, err := host.GetHostByID(dbConnection, h.JumphostID.UUID.String())
	assert.Nil(t, err)

	assert.Equal(t, uuid.MustParse("7f1620cb-b627-4f32-a7eb-e05d19d6065c"), j.IdentityID)
	assert.Equal(t, "test", j.User)
	assert.Equal(t, "jumpbox.com", j.Address)
	assert.Equal(t, 203, j.Port)

	assert.Equal(t, j.Id, h.JumphostID.UUID)

	assert.NotContains(t, h.ExtraFlags, "-v")
	assert.Contains(t, h.ExtraFlags, "-4")

}

func TestCommandToHostWithAlias(t *testing.T) {
	dbConnection := db.GetTestDB(t)
	defer dbConnection.Close()

	id := MockIdentityData()

	err := id.Store(dbConnection)
	assert.Nil(t, err)
	command := `alias ssh_evm='ssh testuser@aliash-test.com -i /Users/test/.ssh/test_key -v'`

	if aliasRe.MatchString(command) {
		match := aliasExtractRe.FindStringSubmatch(command)
		if len(match) > 1 {
			command = match[1]
		}
	}

	h, err := commandToHost(command, dbConnection)
	assert.Nil(t, err)
	assert.Equal(t, "aliash-test.com", h.Address)
	assert.Equal(t, 22, h.Port)
	assert.Equal(t, "testuser", h.User)
	assert.Equal(t, uuid.MustParse("7f1620cb-b627-4f32-a7eb-e05d19d6065c"), h.IdentityID)
}
