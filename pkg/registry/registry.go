package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	natsConfig "github.com/edgefarm/anck/pkg/nats"
	api "github.com/edgefarm/nats-leafnode-sidecar/pkg/api"
	"github.com/edgefarm/nats-leafnode-sidecar/pkg/common"
	"github.com/nats-io/nats.go"
)

const (
	connectTimeoutSeconds = 10
	ngsHost               = "tls://connect.ngs.global:7422"
)

// Registry is a registry for nats-leafnodes
type Registry struct {
	// configFileContent      string
	credsFilesPath         string
	configFilePath         string
	natsConn               *nats.Conn
	registerSubscription   *nats.Subscription
	unregisterSubscription *nats.Subscription
	state                  *State
	config                 *natsConfig.Config
}

// NewRegistry creates a new registry
func NewRegistry(natsConfigPath string, creds string, natsURI string, state string) (*Registry, error) {
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
	config, err := natsConfig.LoadFromFile(natsConfigPath)
	if err != nil {
		return nil, err
	}

	r := &Registry{
		credsFilesPath: creds,
		configFilePath: natsConfigPath,
		natsConn:       nc,
		state:          NewState(state),
		config:         config,
	}
	err = r.updateConfigFile()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Start starts the registry and handles all incoming requests for registering and unregistering
func (r *Registry) Start() error {
	var err error
	r.registerSubscription, err = r.natsConn.Subscribe(common.RegisterSubject, func(m *nats.Msg) {
		creds := &api.Credentials{}
		err := json.Unmarshal(m.Data, creds)
		if err != nil {
			log.Println("Error unmarshalling credentials: ", err)
		}
		fmt.Printf("Received register request for network: %s and component: %s\n", creds.Network, creds.Component)
		err = r.addCredentials(creds.Network, creds.Component, creds.AccountPublicKey)
		if err == nil {
			err = r.natsConn.Publish(m.Reply, []byte(common.OkResponse))
			if err != nil {
				log.Println(err)
			}
		} else {
			err = r.natsConn.Publish(m.Reply, []byte(fmt.Sprintf("%s: %s", common.ErrorResponse, err)))
			if err != nil {
				log.Println(err)
			}
		}
		err = r.writeFile(r.credsFile(creds.Network), creds.Creds)
		if err != nil {
			log.Println(err)
		}
		err = r.updateConfigFile()
		if err != nil {
			log.Println(err)
		}
	})
	if err != nil {
		return err
	}

	r.unregisterSubscription, err = r.natsConn.Subscribe(common.UnregisterSubject, func(m *nats.Msg) {
		log.Println("Received unregister request")
		creds := &api.Credentials{}
		err := json.Unmarshal(m.Data, creds)
		if err != nil {
			log.Println("Error unmarshalling credentials: ", err)
		}
		deleteCredsfile, err := r.removeCredentials(creds.Network, creds.Component)
		if err == nil {
			err = r.natsConn.Publish(m.Reply, []byte(common.OkResponse))
			if err != nil {
				log.Println(err)
			}
		} else {
			err = r.natsConn.Publish(m.Reply, []byte(fmt.Sprintf("%s: %s", common.ErrorResponse, err)))
			if err != nil {
				log.Println(err)
			}
		}

		if deleteCredsfile {
			err = r.removeFile(r.credsFile(creds.Network))
			if err != nil {
				log.Println(err)
			}
			err = r.updateConfigFile()
			if err != nil {
				log.Println(err)
			}
		}
	})
	if err != nil {
		return err
	}
	return nil
}

// Shutdown shuts down the registry
func (r *Registry) Shutdown() {
	log.Println("Shutting down registry")
	if r.registerSubscription != nil {
		r.registerSubscription.Unsubscribe()
	}
	if r.unregisterSubscription != nil {
		r.unregisterSubscription.Unsubscribe()
	}
	r.natsConn.Close()
	os.Exit(0)
}

func (r *Registry) addCredentials(network string, component string, accountPublicKey string) error {
	found := false
	for _, remote := range r.config.Leafnodes.Remotes {
		networkName := filepath.Base(remote.Credentials)
		if networkName == fmt.Sprintf("%s.creds", network) {
			found = true
			break
		}
	}
	if !found {
		err := r.config.AddRemote(ngsHost, r.credsFile(network), accountPublicKey)
		if err != nil {
			return err
		}
	}
	err := r.state.Update(network, component, Add)
	if err != nil {
		return err
	}
	return nil
}

func (r *Registry) removeCredentials(network string, component string) (bool, error) {
	usage, err := r.state.Usage(network)
	if err != nil {
		return false, err
	}
	if usage > 0 {
		fmt.Printf("Removing participant cound from network '%s'\n", network)
		err = r.state.Update(network, component, Remove)
		if err != nil {
			return false, err
		}
		usage--
	}
	if usage <= 0 {
		for _, remote := range r.config.Leafnodes.Remotes {
			if remote.Credentials == r.credsFile(network) {
				err = r.config.RemoveRemoteByCredsfile(remote.Credentials)
				if err != nil {
					return false, err
				}
				break
			}
		}
		err = r.state.Delete(network)
		if err != nil {
			return false, err
		}
	}
	return usage <= 0, nil
}

func (r *Registry) credsFile(network string) string {
	return fmt.Sprintf("%s/%s.creds", r.credsFilesPath, network)
}
