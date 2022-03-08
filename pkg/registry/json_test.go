package registry

import (
	"io/ioutil"
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/stretchr/testify/assert"
)

func TestNewJsonNotExists(t *testing.T) {
	assert := assert.New(t)
	// random, err := ioutil.TempFile("", "")
	// assert.Nil(err)
	config := NewJson("/this/path/does/not/exist")
	config.Dump()
	assert.Nil(nil)
}

func TestNewJsonNotValid(t *testing.T) {
	assert := assert.New(t)
	random, err := ioutil.TempFile("", "")
	assert.Nil(err)

	// No need to check whether `recover()` is nil. Just turn off the panic.
	defer func() { recover() }()

	// this should panic, because the file is not valid json
	NewJson(random.Name())

	// Never reaches here if `NewJson` panics.
	t.Errorf("did not panic")
}

func TestNewJsonValidPresent(t *testing.T) {
	assert := assert.New(t)
	random, err := ioutil.TempFile("", "")
	assert.Nil(err)

	preConfig := `{
	"pid_file": "/var/run/nats/nats.pid",
	"http": 8222,
	"leafnodes": {
	  "remotes": [
		{
		  "url": "nats://localhost:8222",
		  "credentials": "/path/to/creds1.creds"
		},
		{
		  "url": "nats://localhost:8212",
		  "credentials": "/path/to/creds2.creds"
		}
	  ]
	}
}`
	_, err = random.Write([]byte(preConfig))
	assert.Nil(err)

	config := NewJson(random.Name())

	assert.Equal(config.Http, 8222)
	assert.Equal(config.PidFile, "/var/run/nats/nats.pid")
	assert.Equal(config.Leafnodes.Remotes[0].Url, "nats://localhost:8222")
	assert.Equal(config.Leafnodes.Remotes[0].Credentials, "/path/to/creds1.creds")
	assert.Equal(config.Leafnodes.Remotes[1].Url, "nats://localhost:8212")
	assert.Equal(config.Leafnodes.Remotes[1].Credentials, "/path/to/creds2.creds")

	readback := jsonPrettyPrint(config.Dump())
	assert.True(Equal(preConfig, readback))
}

// func TestCredentialsHandling(t *testing.T) {
// 	assert := assert.New(t)
// 	r := &Registry{
// 		configFileContent: defaultConfig,
// 		credsFilesPath:    "/creds",
// 		configFilePath:    "",
// 		natsConn:          nil,
// 	}
// 	err := r.addCredentials("account1", "account1-user")
// 	assert.Nil(err)
// 	r.Dump()
// 	assert.True(Equal(r.Config(), `{
// 		"accounts": {
// 			"account1": {
// 				"users": [{
// 					"password": "account1-password",
// 					"user": "account1-user"
// 				}]
// 			},
// 			"default": {
// 				"users": [{
// 					"user": "default",
// 					"password": ""
// 				}]
// 			}
// 		},
// 		"http": 8222,
// 		"leafnodes": {
// 			"remotes": [
// 				{
// 					"url": "tls://connect.ngs.global:7422",
// 					"credentials": "/creds/account1.creds",
// 					"account": "account1"
// 				}
// 			]
// 		},
// 		"pid_file": "/var/run/nats/nats.pid"
// 	}`))

// 	err = r.addCredentials("account2", "account2-user")
// 	assert.Nil(err)
// 	r.Dump()
// 	assert.True(Equal(r.Config(), `{
// 		"accounts": {
// 			"account1": {
// 				"users": [{
// 					"password": "account1-password",
// 					"user": "account1-user"
// 				}]
// 			},
// 			"account2": {
// 				"users": [{
// 					"password": "account2-password",
// 					"user": "account2-user"
// 				}]
// 			},
// 			"default": {
// 				"users": [{
// 					"user": "default",
// 					"password": ""
// 				}]
// 			}
// 		},
// 		"http": 8222,
// 		"leafnodes": {
// 			"remotes": [
// 				{
// 					"url": "tls://connect.ngs.global:7422",
// 					"credentials": "/creds/account1.creds",
// 					"account": "account1"
// 				},
// 				{
// 					"url": "tls://connect.ngs.global:7422",
// 					"credentials": "/creds/account2.creds",
// 					"account": "account2"
// 				}
// 			]
// 		},
// 		"pid_file": "/var/run/nats/nats.pid"
// 	}`))
// 	err = r.removeCredentials("account1")
// 	assert.Nil(err)
// 	r.Dump()
// 	assert.True(Equal(r.Config(), `{
// 		"accounts": {
// 			"account2": {
// 				"users": [{
// 					"password": "account2-password",
// 					"user": "account2-user"
// 				}]
// 			},
// 			"default": {
// 				"users": [{
// 					"user": "default",
// 					"password": ""
// 				}]
// 			}
// 		},
// 		"http": 8222,
// 		"leafnodes": {
// 			"remotes": [
// 				{
// 					"account": "account2",
// 					"credentials": "/creds/account2.creds",
// 					"url": "tls://connect.ngs.global:7422"
// 				}
// 			]
// 		},
// 		"pid_file": "/var/run/nats/nats.pid"
// 	}`))

// 	// try to remove non existent account
// 	err = r.removeCredentials("account3")
// 	assert.NotNil(err)
// 	r.Dump()
// 	assert.True(Equal(r.Config(), `{
// 		"accounts": {
// 			"account2": {
// 				"users": [{
// 					"password": "account2-password",
// 					"user": "account2-user"
// 				}]
// 			},
// 			"default": {
// 				"users": [{
// 					"user": "default",
// 					"password": ""
// 				}]
// 			}
// 		},
// 		"http": 8222,
// 		"leafnodes": {
// 			"remotes": [
// 				{
// 					"account": "account2",
// 					"credentials": "/creds/account2.creds",
// 					"url": "tls://connect.ngs.global:7422"
// 				}
// 			]
// 		},
// 		"pid_file": "/var/run/nats/nats.pid"
// 	}`))
// 	err = r.removeCredentials("account2")
// 	assert.Nil(err)
// 	r.Dump()
// 	assert.True(Equal(r.Config(), `{
// 		"accounts": {
// 			"default": {
// 				"users": [{
// 					"user": "default",
// 					"password": ""
// 				}]
// 			}},
// 		"http": 8222,
// 		"leafnodes": {
// 			"remotes": []
// 		},
// 		"pid_file": "/var/run/nats/nats.pid"
// 	}`))
// 	err = r.addCredentials("account3", "account3-user")
// 	assert.Nil(err)
// 	r.Dump()
// 	assert.True(Equal(r.Config(), `{
// 		"accounts": {
// 			"account3": {
// 				"users": [{
// 					"password": "account3-password",
// 					"user": "account3-user"
// 				}]
// 			},
// 			"default": {
// 				"users": [{
// 					"user": "default",
// 					"password": ""
// 				}]
// 			}
// 		},
// 		"http": 8222,
// 		"leafnodes": {
// 			"remotes": [
// 				{
// 					"url": "tls://connect.ngs.global:7422",
// 					"credentials": "/creds/account3.creds",
// 					"account": "account3"
// 				}
// 			]
// 		},
// 		"pid_file": "/var/run/nats/nats.pid"
// 	}`))
// 	err = r.addCredentials("account1", "account1-user")
// 	assert.Nil(err)
// 	r.Dump()
// 	assert.True(Equal(r.Config(), `{
// 		"accounts": {
// 			"account1": {
// 				"users": [{
// 					"password": "account1-password",
// 					"user": "account1-user"
// 				}]
// 			},
// 			"account3": {
// 				"users": [{
// 					"password": "account3-password",
// 					"user": "account3-user"
// 				}]
// 			},
// 			"default": {
// 				"users": [{
// 					"user": "default",
// 					"password": ""
// 				}]
// 			}
// 		},
// 		"http": 8222,
// 		"leafnodes": {
// 			"remotes": [
// 				{
// 					"url": "tls://connect.ngs.global:7422",
// 					"credentials": "/creds/account3.creds",
// 					"account": "account3"
// 				},
// 				{
// 					"url": "tls://connect.ngs.global:7422",
// 					"credentials": "/creds/account1.creds",
// 					"account": "account1"
// 				}
// 			]
// 		},
// 		"pid_file": "/var/run/nats/nats.pid"
// 	}`))
// 	// update credentials
// 	err = r.addCredentials("account1", "account1-user")
// 	assert.Nil(err)
// 	r.Dump()
// 	assert.True(Equal(r.Config(), `{
// 		"accounts": {
// 			"account1": {
// 				"users": [{
// 					"password": "account1-newpassword",
// 					"user": "account1-user"
// 				}]
// 			},
// 			"account3": {
// 				"users": [{
// 					"password": "account3-password",
// 					"user": "account3-user"
// 				}]
// 			},
// 			"default": {
// 				"users": [{
// 					"user": "default",
// 					"password": ""
// 				}]
// 			}
// 		},
// 		"http": 8222,
// 		"leafnodes": {
// 			"remotes": [
// 				{
// 					"url": "tls://connect.ngs.global:7422",
// 					"credentials": "/creds/account3.creds",
// 					"account": "account3"
// 				},
// 				{
// 					"url": "tls://connect.ngs.global:7422",
// 					"credentials": "/creds/account1.creds",
// 					"account": "account1"
// 				}
// 			]
// 		},
// 		"pid_file": "/var/run/nats/nats.pid"
// 	}`))
// }

func Equal(configFileContent string, reference string) bool {
	return jsonpatch.Equal([]byte(reference), []byte(configFileContent))
}
