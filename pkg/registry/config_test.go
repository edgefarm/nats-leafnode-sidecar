package registry

// import (
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestUpdateConfigFile(t *testing.T) {
// 	assert := assert.New(t)
// 	r := &Registry{
// 		configFileContent: `{
// 			"accounts": {},
// 			"http": 8222,
// 			"leafnodes": {
// 				"remotes": []
// 			},
// 			"pid_file": "/var/run/nats.pid",
// 			"server_name": "edge"
// 		}`,
// 		credsFilesPath: "/tmp/test/creds",
// 		configFilePath: "/tmp/test/config",
// 		natsConn:       nil,
// 	}
// 	file, err := os.CreateTemp("", "TestUpdateConfigFile")
// 	assert.Nil(err)
// 	r.configFilePath = file.Name()
// 	defer os.Remove(r.configFilePath)
// 	err = r.addCredentials("account1", "account1-user")
// 	assert.Nil(err)
// 	err = r.updateConfigFile()
// 	assert.Nil(err)
// 	dat, err := os.ReadFile(r.configFilePath)
// 	assert.Nil(err)
// 	assert.True(Equal(r.Config(), string(dat)))
// }

// func TestUpdateConfigFileWithSymlink(t *testing.T) {
// 	assert := assert.New(t)
// 	r := &Registry{
// 		configFileContent: `{
// 			"accounts": {},
// 			"http": 8222,
// 			"leafnodes": {
// 				"remotes": []
// 			},
// 			"pid_file": "/var/run/nats.pid",
// 			"server_name": "edge"
// 		}`,
// 		credsFilesPath: "/tmp/test/creds",
// 		configFilePath: "/tmp/test/config",
// 		natsConn:       nil,
// 	}
// 	file, err := os.CreateTemp("", "TestUpdateConfigFileWithSymlink")
// 	assert.Nil(err)
// 	symlink := fmt.Sprintf("%s.link", file.Name())
// 	err = os.Symlink(file.Name(), symlink)
// 	assert.Nil(err)
// 	r.configFilePath = symlink
// 	defer os.Remove(file.Name())
// 	defer os.Remove(symlink)
// 	err = r.addCredentials("account1", "account1-user")
// 	assert.Nil(err)
// 	err = r.updateConfigFile()
// 	assert.Nil(err)
// 	dat, err := os.ReadFile(r.configFilePath)
// 	assert.Nil(err)
// 	assert.True(Equal(r.Config(), string(dat)))
// }
