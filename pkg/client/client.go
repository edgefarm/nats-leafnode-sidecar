/*
Copyright © 2021 Ci4Rail GmbH <engineering@ci4rail.com>
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	files "github.com/edgefarm/nats-leafnode-sidecar/pkg/files"
	nats "github.com/nats-io/nats.go"
)

const (
	credentialsMountDirectory        = "/nats-credentials"
	edgefarmNetworkAccountNameSecret = "edgefarm.network-natsAccount"
)

// NatsCredentials contains the credentials for the nats server.
type NatsCredentials struct {
	Username         string `json:"username"`
	CredsFileContent string `json:"creds"`
}

// Client is a client for the registry service.
type Client struct {
	// NatsAccount is the nats account to use for the given nats user.
	NatsAccount string
	// NatsUser is the user name for the nats server for the given nats account.
	NatsUser string
	// Creds contains the credentials for the current application
	Creds *NatsCredentials
	// Nats connection
	NatsConn *nats.Conn
}

// NewClient creates a new client for the registry service.
func NewClient(natsUser string, natsURI string) (*Client, error) {
	f, err := files.GetSymlinks(credentialsMountDirectory)
	if err != nil {
		return nil, err
	}

	natsAccount, err := func() (string, error) {
		f, err := files.GetSymlinks(credentialsMountDirectory)
		if err != nil {
			return "", err
		}
		for _, file := range f {
			if filepath.Base(file) == edgefarmNetworkAccountNameSecret {
				b, err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Println(err)
					break
				}
				return string(b), nil
			}
		}
		return "", fmt.Errorf("no nats account found for user %s", natsUser)
	}()
	if err != nil {
		return nil, err
	}

	creds := func() *NatsCredentials {
		for _, file := range f {
			isDir, err := files.IsDir(file)
			if err != nil {
				break
			}
			if !isDir {
				b, err := ioutil.ReadFile(file)
				if err != nil {
					fmt.Println(err)
					break
				}
				if natsUser == filepath.Base(file) {
					return &NatsCredentials{
						Username:         fmt.Sprintf("%s-%s", natsAccount, filepath.Base(file)),
						CredsFileContent: string(b),
					}
				}
			}
		}
		return nil
	}()

	if creds == nil {
		return nil, fmt.Errorf("no credentials found for user %s", natsUser)
	}
	// nc, err := nats.Connect(natsURI)
	// if err != nil {
	// 	return nil, err
	// }

	return &Client{
		NatsAccount: natsAccount,
		NatsUser:    natsUser,
		Creds:       creds,
		NatsConn:    nil,
	}, nil
}

// Connect registeres the application and connects to the nats server.
func (c *Client) Connect() error {
	fmt.Printf("%+v", c.Creds)
	err := c.Registry(Register())
	if err != nil {
		return err
	}

	return nil
}

// Shutdown unregisteres the application and shuts down the nats connection.
func (c *Client) Shutdown() error {
	err := c.Registry(Unregister())
	if err != nil {
		return err
	}
	// c.NatsConn.Close()
	return nil
}
