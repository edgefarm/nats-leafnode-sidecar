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
	"log"
	"path/filepath"
	"time"

	api "github.com/edgefarm/edgefarm.network/pkg/apis/config/v1alpha1"
	common "github.com/edgefarm/nats-leafnode-sidecar/pkg/common"
	files "github.com/edgefarm/nats-leafnode-sidecar/pkg/files"
	nats "github.com/nats-io/nats.go"
)

const (
	edgefarmNetworkAccountNameSecret = "edgefarm.network-natsUserData"
	connectTimeoutSeconds            = 10
)

// NatsCredentials contains the credentials for the nats server.
type NatsCredentials struct {
	Username         string `json:"username"`
	CredsFileContent string `json:"creds"`
}

// Client is a client for the registry service.
type Client struct {
	// creds contains the credentials for the current application
	creds *api.Credentials
	// Nats connection
	nc *nats.Conn
}

// NewClient creates a new client for the registry service.
func NewClient(credentialsMountDirectory string, natsURI string) (*Client, error) {
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

	opts := []nats.Option{nats.Timeout(time.Duration(1) * time.Second)}
	opts = common.SetupConnOptions(opts)
	ncChan := make(chan *nats.Conn)
	go func() {
		for {
			log.Printf("\rConnecting to nats server: %s\n", natsURI)
			nc, err := nats.Connect(natsURI, opts...)
			if err != nil {
				log.Printf("Connect failed to %s: %s\n", natsURI, err)
			} else {
				log.Printf("Connected to '%s'\n", natsURI)
				ncChan <- nc
				return
			}
			func() {
				for i := connectTimeoutSeconds; i >= 0; i-- {
					time.Sleep(time.Second)
					log.Printf("\rReconnecting in %2d seconds", i)
				}
				log.Println("")
			}()
		}
	}()

	nc := <-ncChan

	return &Client{
		creds: creds,
		nc:    nc,
	}, nil
}

// Connect registeres the application and connects to the nats server.
func (c *Client) Connect() error {
	log.Printf("Credentials found for userAccountName %s\n", c.creds.UserAccountName)
	err := c.Registry(Register())
	if err != nil {
		return err
	}

	return nil
}

// Shutdown unregisteres the application and shuts down the nats connection.
func (c *Client) Shutdown() error {
	log.Println("Shutting down client")
	err := c.Registry(Unregister())
	if err != nil {
		return err
	}
	c.nc.Close()
	return nil
}
