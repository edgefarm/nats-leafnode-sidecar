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
	"path/filepath"
	"strings"
	"time"

	api "github.com/edgefarm/nats-leafnode-sidecar/pkg/api"
	common "github.com/edgefarm/nats-leafnode-sidecar/pkg/common"
	"github.com/edgefarm/nats-leafnode-sidecar/pkg/files"
	"github.com/fsnotify/fsnotify"
	nats "github.com/nats-io/nats.go"
)

const (
	connectTimeoutSeconds = 10
)

var (
	// IgnoredFromWatch is a list of files that should be ignored from the watch.
	ignoredFromWatch = []string{"edgefarm-sys.creds", "..data", ".pub"}
)

// NatsCredentials contains the credentials for the nats server.
type NatsCredentials struct {
	Username         string `json:"username"`
	CredsFileContent string `json:"creds"`
}

// Client is a client for the registry service.
type Client struct {
	// component is the name of the component this client is for.
	component string
	// path to the credentials files
	path string
	// Nats connection
	nc *nats.Conn
	// Watcher to monitor credentials directory
	watcher *fsnotify.Watcher
	// finish is a channel to signal the client to shutdown.
	finish chan interface{}
	// finishWatch is a channel to signal the watch loop to finish
	finishWatch chan interface{}
}

// NewClient creates a new client for the registry service.
func NewClient(credentialsMountDirectory string, natsURI string, component string) (*Client, error) {
	opts := []nats.Option{nats.Timeout(time.Duration(1) * time.Second)}
	opts = append(opts, nats.UserCredentials(common.CredentialsFile))
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
		component:   component,
		path:        credentialsMountDirectory,
		nc:          nc,
		watcher:     watcher,
		finish:      make(chan interface{}),
		finishWatch: make(chan interface{}),
	}, nil
}

// Start starts the client.
func (c *Client) Start() error {
	go c.loop()
	return nil
}

func (c *Client) add() error {
	return c.action(Register())
}

func (c *Client) removeAll() error {
	return c.action(Unregister())
}

func isIgnored(file string) bool {
	for _, i := range ignoredFromWatch {
		if strings.Contains(file, i) {
			return true
		}
	}
	return false
}

func (c *Client) action(option *RegistryOptions) error {
	f, err := files.GetSymlinks(c.path)
	if err != nil {
		return err
	}
	credsFiles := make([]string, 0)
	for _, f := range f {
		if !isIgnored(f) {
			credsFiles = append(credsFiles, f)
		}
	}

	for _, file := range credsFiles {
		credsContent, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		accountPubKeyContent, err := ioutil.ReadFile(fmt.Sprintf("%s.pub", file))
		if err != nil {
			return err
		}
		networkName := filepath.Base(file)
		creds := &api.Credentials{
			Network:          filepath.Base(file),
			Component:        c.component,
			Creds:            string(credsContent),
			AccountPublicKey: string(accountPubKeyContent),
		}
		fmt.Printf("%s network %s\n", option.subject, networkName)
		err = c.Registry(option, creds)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) remove(networkPath string) error {
	network := filepath.Base(networkPath)
	network = strings.TrimSuffix(network, ".creds")
	creds := &api.Credentials{
		Network:   network,
		Component: c.component,
	}
	fmt.Printf("Unregistering network %s\n", network)
	err := c.Registry(Unregister(), creds)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) installWatch(path string, addCallback func() error, removeCallback func(string) error) error {
	go func() {
		for {
			select {
			case event, ok := <-c.watcher.Events:
				if !ok {
					return
				}
				ignored := false
				for _, ignoredFile := range ignoredFromWatch {
					if strings.Contains(event.Name, ignoredFile) {
						ignored = true
					}
				}
				if ignored {
					fmt.Println("Ignoring event: ", event)
					continue
				}
				fmt.Println("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created/modified file:", event.Name)
					err := addCallback()
					if err != nil {
						log.Println(err)
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("removed/renamed file:", event.Name)
					err := removeCallback(event.Name)
					if err != nil {
						log.Println(err)
					}
				}
			case err, ok := <-c.watcher.Errors:
				fmt.Println("c.Watcher.Errors: ", err, ok)
				if !ok {
					return
				}
				log.Println("error:", err)

			case <-c.finishWatch:
				fmt.Println("Stopping watcher")
				c.watcher.Close()
				return

			case <-time.After(1 * time.Second):
			}
		}
	}()

	err := c.watcher.Add(path)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) addCallback() error {
	return c.add()
}

func (c *Client) removeCallback(network string) error {
	return c.remove(network)
}

// loop runs the client in a loop.
func (c *Client) loop() {
	// first time register all the credentials
	err := c.add()
	if err != nil {
		log.Println(err)
	}

	// the watch will re-register the credentials on changes
	err = c.installWatch(c.path, c.addCallback, c.removeCallback)
	if err != nil {
		log.Println(err)
	}

	for {
		select {
		case <-c.finish:
			fmt.Println("Stopping loop")
			return
		default:
			time.Sleep(time.Second * 1)
		}
	}
}

// Shutdown unregisteres the application and shuts down the nats connection.
func (c *Client) Shutdown() {
	log.Println("Shutting down client")
	err := c.removeAll()
	if err != nil {
		fmt.Println(err)
	}
	c.finishWatch <- true
	c.finish <- true
	c.nc.Close()
}
