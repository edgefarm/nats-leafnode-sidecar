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
	"log"
	"time"

	api "github.com/edgefarm/anck-credentials/pkg/apis/config/v1alpha1"
	common "github.com/edgefarm/nats-leafnode-sidecar/pkg/common"
	"github.com/edgefarm/nats-leafnode-sidecar/pkg/files"
	"github.com/fsnotify/fsnotify"
	nats "github.com/nats-io/nats.go"
)

const (
	connectTimeoutSeconds = 10
)

// NatsCredentials contains the credentials for the nats server.
type NatsCredentials struct {
	Username         string `json:"username"`
	CredsFileContent string `json:"creds"`
}

// Client is a client for the registry service.
type Client struct {
	// path to the credentials files
	path string
	// Nats connection
	nc *nats.Conn
	// Watcher to monitor credentials directory
	watcher *fsnotify.Watcher
	// Reregister
	reregister chan interface{}
	// finish
	finish chan interface{}
}

// NewClient creates a new client for the registry service.
func NewClient(credentialsMountDirectory string, natsURI string) (*Client, error) {
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

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Client{
		path:       credentialsMountDirectory,
		nc:         nc,
		watcher:    watcher,
		reregister: make(chan interface{}),
		finish:     make(chan interface{}),
	}, nil
}

func (c *Client) Action(option *RegistryOptions) error {
	f, err := files.GetSymlinks(c.path)
	if err != nil {
		return err
	}
	for _, file := range f {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		creds := &api.Credentials{
			NetworkParticipant: file,
			Creds:              string(b),
		}
		fmt.Printf("%s network %s\n", option.subject, file)
		err = c.Registry(option, creds)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) Watch(path string, callback func(string)) error {
	go func() {
		for {
			select {
			case event, ok := <-c.watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					callback(event.Name)
				}
			case err, ok := <-c.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err := c.watcher.Add(path)
	if err != nil {
		return err
	}
	return nil
}

// Loop runs the client in a loop.
func (c *Client) Loop() {
	err := c.Action(Register())
	if err != nil {
		log.Println(err)
	}

	go Watch(c.path, c.reregister)
	for {
		log.Println("Looping")
		time.Sleep(time.Second * 5)
	}
}

// // Connect registeres the application and connects to the nats server.
// func (c *Client) Connect() error {
// 	log.Printf("Credentials found for userAccountName %s\n", c.creds.Creds)
// 	err := c.Registry(Register())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// Shutdown unregisteres the application and shuts down the nats connection.
func (c *Client) Shutdown() error {
	log.Println("Shutting down client")
	err := c.Registry(Unregister())
	if err != nil {
		return err
	}
	c.nc.Close()
	c.watcher.Close()
	return nil
}
