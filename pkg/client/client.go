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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	api "github.com/edgefarm/edgefarm.network/pkg/apis/config/v1alpha1"
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
	// Creds contains the credentials for the current application
	Creds *api.Credentials
	// Nats connection
	NatsConn *nats.Conn
}

// NewClient creates a new client for the registry service.
func NewClient(natsURI string) (*Client, error) {
	creds := &api.Credentials{}
	err := func() error {
		f, err := files.GetSymlinks(credentialsMountDirectory)
		if err != nil {
			return err
		}
		for _, file := range f {
			if filepath.Base(file) == edgefarmNetworkAccountNameSecret {
				b, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}
				err = json.Unmarshal(b, &creds)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return fmt.Errorf("no credentials file found at '%s/%s'", credentialsMountDirectory, edgefarmNetworkAccountNameSecret)
	}()
	if err != nil {
		return nil, err
	}

	nc, err := nats.Connect(natsURI)
	if err != nil {
		return nil, err
	}

	return &Client{
		Creds:    creds,
		NatsConn: nc,
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
	fmt.Println("Shutting down client")
	err := c.Registry(Unregister())
	if err != nil {
		return err
	}
	c.NatsConn.Close()
	return nil
}
