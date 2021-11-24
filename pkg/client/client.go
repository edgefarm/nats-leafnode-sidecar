package client

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	files "github.com/edgefarm/nats-leafnode-sidecar/pkg/files"
	nats "github.com/nats-io/nats.go"
)

const (
	credentialsMountDirectory = "/nats-credentials"
)

type NatsCredentials struct {
	Username         string `json:"username"`
	CredsFileContent string `json:"creds"`
}

// Client is a client for the registry service.
type Client struct {
	// NatsAccount is the account name for the nats server.
	NatsAccount string
	// Creds contains the credentials for the current application
	Creds []NatsCredentials
	// Nats connection
	NatsConn *nats.Conn
}

func NewClient(natsAccount string, natsURI string) (*Client, error) {
	f, err := files.GetSymlinks(credentialsMountDirectory)
	if err != nil {
		return nil, err
	}
	creds := func() []NatsCredentials {
		var creds []NatsCredentials
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
				creds = append(creds, NatsCredentials{
					Username:         filepath.Base(file),
					CredsFileContent: string(b),
				})
			}
		}
		return creds
	}()

	// nc, err := nats.Connect(natsURI)
	// if err != nil {
	// 	return nil, err
	// }

	return &Client{
		NatsAccount: natsAccount,
		Creds:       creds,
		NatsConn:    nil,
	}, nil
}

func (c *Client) Connect() error {
	fmt.Printf("%+v", c.Creds)
	err := c.Registry(Register())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Shutdown() error {
	err := c.Registry(Unregister())
	if err != nil {
		return err
	}
	// c.NatsConn.Close()
	return nil
}
