package host

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConnectionString(t *testing.T) {
	err := ValidateConnectionString("ssh -i ${identitiy_file_path} -p ${port} ${user}:${address}")
	assert.Nil(t, err)

	err = ValidateConnectionString("ssh -i ${identitiy} -p ${port} ${user}:${address}")
	assert.Error(t, err)
}
