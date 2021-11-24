package client

import (
	"encoding/json"
)

const (
	registerSubject   = "register"
	unregisterSubject = "unregister"
)

type RegistryOptions struct {
	subject string
}

func Register() *RegistryOptions {
	return &RegistryOptions{
		subject: registerSubject,
	}
}

func Unregister() *RegistryOptions {
	return &RegistryOptions{
		subject: unregisterSubject,
	}
}

func (c *Client) Registry(option *RegistryOptions) error {
	secret := []byte{}
	for _, cred := range c.Creds {
		if cred.Username == c.NatsAccount {
			var err error
			secret, err = json.Marshal(cred)
			if err != nil {
				return err
			}
		}
	}
	if len(secret) > 0 {
		return c.NatsConn.Publish(option.subject, secret)
	}
	return nil
}
