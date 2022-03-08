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
	"time"

	api "github.com/edgefarm/anck-credentials/pkg/apis/config/v1alpha1"
	"github.com/edgefarm/nats-leafnode-sidecar/pkg/common"
)

const (
	natsTimeout = time.Second * 10
)

// RegistryOptions is used to configure a Registry request
type RegistryOptions struct {
	subject string
}

// Register is used to configure a Register request
func Register() *RegistryOptions {
	return &RegistryOptions{
		subject: common.RegisterSubject,
	}
}

// Unregister is used to configure an Unregister request
func Unregister() *RegistryOptions {
	return &RegistryOptions{
		subject: common.UnregisterSubject,
	}
}

// Registry is used to register or unregister an application to the nats server
func (c *Client) Registry(option *RegistryOptions, creds *api.Credentials) error {
	j, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	resp, err := c.nc.Request(option.subject, j, natsTimeout)
	if err != nil {
		return err
	}
	if string(resp.Data) != common.OkResponse {
		return fmt.Errorf("request failed: %s", string(resp.Data))
	}
	return nil
}
