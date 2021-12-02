package registry

import (
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/stretchr/testify/assert"
)

func TestCredentialsHandling(t *testing.T) {
	assert := assert.New(t)
	r := &Registry{
		configFileContent: `{
			"accounts": {},
			"http": 8222,
			"leafnodes": {
				"remotes": []
			},
			"pid_file": "/var/run/nats.pid",
			"server_name": "edge"
		}`,
		credsFilesPath: "",
		configFilePath: "",
		natsConn:       nil,
	}
	err := r.addCredentials("account1", "account1-user", "account1-password", "/account1-user.creds")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account1": {
				"users": {
					"password": "account1-password",
					"user": "account1-user"
				}
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "account1-user.creds",
					"account": "account1"
				}
			]
		},
		"pid_file": "/var/run/nats.pid",
		"server_name": "edge"
	}`))
	err = r.addCredentials("account2", "account2-user", "account2-password", "/account2-user.creds")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account1": {
				"users": {
					"password": "account1-password",
					"user": "account1-user"
				}
			},
			"account2": {
				"users": {
					"password": "account2-password",
					"user": "account2-user"
				}
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "account1-user.creds",
					"account": "account1"
				},
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "account2-user.creds",
					"account": "account2"
				}
			]
		},
		"pid_file": "/var/run/nats.pid",
		"server_name": "edge"
	}`))
	err = r.removeCredentials("account1")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account2": {
				"users": {
					"password": "account2-password",
					"user": "account2-user"
				}
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"account": "account2",
					"credentials": "account2-user.creds",
					"url": "tls://connect.ngs.global:7422"
				}
			]
		},
		"pid_file": "/var/run/nats.pid",
		"server_name": "edge"
	}`))

	// try to remove non existent account
	err = r.removeCredentials("account3")
	assert.NotNil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account2": {
				"users": {
					"password": "account2-password",
					"user": "account2-user"
				}
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"account": "account2",
					"credentials": "account2-user.creds",
					"url": "tls://connect.ngs.global:7422"
				}
			]
		},
		"pid_file": "/var/run/nats.pid",
		"server_name": "edge"
	}`))
	err = r.removeCredentials("account2")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {},
		"http": 8222,
		"leafnodes": {
			"remotes": []
		},
		"pid_file": "/var/run/nats.pid",
		"server_name": "edge"
	}`))
	err = r.addCredentials("account3", "account3-user", "account3-password", "/account3-user.creds")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account3": {
				"users": {
					"password": "account3-password",
					"user": "account3-user"
				}
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "account3-user.creds",
					"account": "account3"
				}
			]
		},
		"pid_file": "/var/run/nats.pid",
		"server_name": "edge"
	}`))
	err = r.addCredentials("account1", "account1-user", "account1-password", "/account1-user.creds")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account1": {
				"users": {
					"password": "account1-password",
					"user": "account1-user"
				}
			},
			"account3": {
				"users": {
					"password": "account3-password",
					"user": "account3-user"
				}
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "account3-user.creds",
					"account": "account3"
				},
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "account1-user.creds",
					"account": "account1"
				}
			]
		},
		"pid_file": "/var/run/nats.pid",
		"server_name": "edge"
	}`))
	// update credentials
	err = r.addCredentials("account1", "account1-user", "account1-newpassword", "/account1-user.creds")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account1": {
				"users": {
					"password": "account1-newpassword",
					"user": "account1-user"
				}
			},
			"account3": {
				"users": {
					"password": "account3-password",
					"user": "account3-user"
				}
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "account3-user.creds",
					"account": "account3"
				},
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "account1-user.creds",
					"account": "account1"
				}
			]
		},
		"pid_file": "/var/run/nats.pid",
		"server_name": "edge"
	}`))
}

func Equal(configFileContent string, reference string) bool {
	return jsonpatch.Equal([]byte(reference), []byte(configFileContent))
}
