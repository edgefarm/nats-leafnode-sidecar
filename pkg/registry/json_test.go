package registry

import (
	"testing"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/stretchr/testify/assert"
)

func TestCredentialsHandling(t *testing.T) {
	assert := assert.New(t)
	r := &Registry{
		configFileContent: config,
		credsFilesPath:    "/creds",
		configFilePath:    "",
		natsConn:          nil,
	}
	err := r.addCredentials("account1", "account1-user", "account1-password")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account1": {
				"users": [{
					"password": "account1-password",
					"user": "account1-user"
				}]
			},
			"default": {
				"users": [{
					"user": "",
					"password": ""
				}]
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "/creds/account1.creds",
					"account": "account1"
				}
			]
		},
		"pid_file": "/var/run/nats/nats.pid"
	}`))

	err = r.addCredentials("account2", "account2-user", "account2-password")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account1": {
				"users": [{
					"password": "account1-password",
					"user": "account1-user"
				}]
			},
			"account2": {
				"users": [{
					"password": "account2-password",
					"user": "account2-user"
				}]
			},
			"default": {
				"users": [{
					"user": "",
					"password": ""
				}]
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "/creds/account1.creds",
					"account": "account1"
				},
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "/creds/account2.creds",
					"account": "account2"
				}
			]
		},
		"pid_file": "/var/run/nats/nats.pid"
	}`))
	err = r.removeCredentials("account1")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account2": {
				"users": [{
					"password": "account2-password",
					"user": "account2-user"
				}]
			},
			"default": {
				"users": [{
					"user": "",
					"password": ""
				}]
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"account": "account2",
					"credentials": "/creds/account2.creds",
					"url": "tls://connect.ngs.global:7422"
				}
			]
		},
		"pid_file": "/var/run/nats/nats.pid"
	}`))

	// try to remove non existent account
	err = r.removeCredentials("account3")
	assert.NotNil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account2": {
				"users": [{
					"password": "account2-password",
					"user": "account2-user"
				}]
			},
			"default": {
				"users": [{
					"user": "",
					"password": ""
				}]
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"account": "account2",
					"credentials": "/creds/account2.creds",
					"url": "tls://connect.ngs.global:7422"
				}
			]
		},
		"pid_file": "/var/run/nats/nats.pid"
	}`))
	err = r.removeCredentials("account2")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"default": {
				"users": [{
					"user": "",
					"password": ""
				}]
			}},
		"http": 8222,
		"leafnodes": {
			"remotes": []
		},
		"pid_file": "/var/run/nats/nats.pid"
	}`))
	err = r.addCredentials("account3", "account3-user", "account3-password")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account3": {
				"users": [{
					"password": "account3-password",
					"user": "account3-user"
				}]
			},
			"default": {
				"users": [{
					"user": "",
					"password": ""
				}]
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "/creds/account3.creds",
					"account": "account3"
				}
			]
		},
		"pid_file": "/var/run/nats/nats.pid"
	}`))
	err = r.addCredentials("account1", "account1-user", "account1-password")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account1": {
				"users": [{
					"password": "account1-password",
					"user": "account1-user"
				}]
			},
			"account3": {
				"users": [{
					"password": "account3-password",
					"user": "account3-user"
				}]
			},
			"default": {
				"users": [{
					"user": "",
					"password": ""
				}]
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "/creds/account3.creds",
					"account": "account3"
				},
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "/creds/account1.creds",
					"account": "account1"
				}
			]
		},
		"pid_file": "/var/run/nats/nats.pid"
	}`))
	// update credentials
	err = r.addCredentials("account1", "account1-user", "account1-newpassword")
	assert.Nil(err)
	r.Dump()
	assert.True(Equal(r.Config(), `{
		"accounts": {
			"account1": {
				"users": [{
					"password": "account1-newpassword",
					"user": "account1-user"
				}]
			},
			"account3": {
				"users": [{
					"password": "account3-password",
					"user": "account3-user"
				}]
			},
			"default": {
				"users": [{
					"user": "",
					"password": ""
				}]
			}
		},
		"http": 8222,
		"leafnodes": {
			"remotes": [
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "/creds/account3.creds",
					"account": "account3"
				},
				{
					"url": "tls://connect.ngs.global:7422",
					"credentials": "/creds/account1.creds",
					"account": "account1"
				}
			]
		},
		"pid_file": "/var/run/nats/nats.pid"
	}`))
}

func Equal(configFileContent string, reference string) bool {
	return jsonpatch.Equal([]byte(reference), []byte(configFileContent))
}
