package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/edgefarm/nats-leafnode-sidecar/pkg/common"
	"github.com/nats-io/nats.go"

	api "github.com/edgefarm/edgefarm.network/pkg/apis/config/v1alpha1"
)

const (
	credsFilesPath = "/tmp/sidecar/creds"
	configFilePath = "/tmp/sidecar/nats.json"
)

// Registry is a registry for nats-leafnodes
type Registry struct {
	configFileContent string
	credsFilesPath    string
	configFilePath    string
	// Nats connection
	natsConn *nats.Conn
}

// NewRegistry creates a new registry
func NewRegistry(natsURI string) (*Registry, error) {
	nc := &nats.Conn{}
	if natsURI != "" {
		var err error
		nc, err = nats.Connect(natsURI)
		if err != nil {
			return nil, err
		}
	}
	r := &Registry{
		configFileContent: string(config),
		credsFilesPath:    credsFilesPath,
		configFilePath:    configFilePath,
		natsConn:          nc,
	}
	return r, nil
}

// Start starts the registry and handles all incoming requests for registering and unregistering
func (r *Registry) Start() error {
	_, err := r.natsConn.Subscribe(common.RegisterSubject, func(m *nats.Msg) {
		fmt.Println("Received register request")
		userCreds := &api.Credentials{}
		err := json.Unmarshal(m.Data, userCreds)
		if err != nil {
			fmt.Println("Error unmarshalling credentials: ", err)
		}
		err = r.addCredentials(userCreds.UserAccountName, userCreds.Username, userCreds.Password, userCreds.Creds)
		if err == nil {
			err = r.natsConn.Publish(m.Reply, []byte(common.OkResponse))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = r.natsConn.Publish(m.Reply, []byte(fmt.Sprintf("%s: %s", common.ErrorResponse, err)))
			if err != nil {
				fmt.Println(err)
			}
		}
		err = r.updateConfigFile()
		if err != nil {
			fmt.Println(err)
		}
		err = r.writeFile(fmt.Sprintf("%s/%s", r.credsFilesPath, userCreds.UserAccountName), userCreds.Creds)
		if err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		return err
	}

	_, err = r.natsConn.Subscribe(common.UnregisterSubject, func(m *nats.Msg) {
		fmt.Println("Received unregister request")
		userCreds := &api.Credentials{}
		err := json.Unmarshal(m.Data, userCreds)
		if err != nil {
			fmt.Println("Error unmarshalling credentials: ", err)
		}
		err = r.removeCredentials(userCreds.UserAccountName)
		if err == nil {
			err = r.natsConn.Publish(m.Reply, []byte(common.OkResponse))
			if err != nil {
				fmt.Println(err)
			}
		} else {
			err = r.natsConn.Publish(m.Reply, []byte(fmt.Sprintf("%s: %s", common.ErrorResponse, err)))
			if err != nil {
				fmt.Println(err)
			}
		}
		err = r.updateConfigFile()
		if err != nil {
			fmt.Println(err)
		}
		err = r.removeFile(fmt.Sprintf("%s/%s", r.credsFilesPath, userCreds.UserAccountName))
		if err != nil {
			fmt.Println(err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}

// Shutdown shuts down the registry
func (r *Registry) Shutdown() {
	fmt.Println("Shutting down registry")
	os.Exit(0)
}
