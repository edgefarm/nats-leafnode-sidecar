package registry

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateConfigFile(t *testing.T) {
	assert := assert.New(t)
	r, err := NewRegistry("")
	assert.Nil(err)
	file, err := os.CreateTemp("", "TestUpdateConfigFile")
	assert.Nil(err)
	r.configFilePath = file.Name()
	defer os.Remove(r.configFilePath)
	err = r.addCredentials("account1", "account1-user", "account1-password", "/account1-user.creds")
	assert.Nil(err)
	err = r.updateConfigFile()
	assert.Nil(err)
	dat, err := os.ReadFile(r.configFilePath)
	assert.Nil(err)
	assert.True(Equal(r.Config(), string(dat)))
}

func TestUpdateConfigFileWithSymlink(t *testing.T) {
	assert := assert.New(t)
	r, err := NewRegistry("")
	assert.Nil(err)
	file, err := os.CreateTemp("", "TestUpdateConfigFileWithSymlink")
	assert.Nil(err)
	symlink := fmt.Sprintf("%s.link", file.Name())
	err = os.Symlink(file.Name(), symlink)
	assert.Nil(err)
	r.configFilePath = symlink
	defer os.Remove(file.Name())
	defer os.Remove(symlink)
	err = r.addCredentials("account1", "account1-user", "account1-password", "/account1-user.creds")
	assert.Nil(err)
	err = r.updateConfigFile()
	assert.Nil(err)
	dat, err := os.ReadFile(r.configFilePath)
	assert.Nil(err)
	assert.True(Equal(r.Config(), string(dat)))
}
