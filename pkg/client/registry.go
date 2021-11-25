package client

import "time"

const (
	natsTimeout       = time.Second * 10
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
	resp, err := c.NatsConn.Request(option.subject, []byte(c.Creds.CredsFileContent), natsTimeout)
	if err != nil {
		return err
	}
	// TODO: check resp.Data
	if resp.Data == nil {
		return nil
	}
	return nil
}
