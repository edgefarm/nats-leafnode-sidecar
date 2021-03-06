package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/edgefarm/anck/pkg/jetstreams"
	natsConfig "github.com/edgefarm/anck/pkg/nats"
	api "github.com/edgefarm/nats-leafnode-sidecar/pkg/api"
	"github.com/edgefarm/nats-leafnode-sidecar/pkg/common"
	"github.com/nats-io/nats.go"
)

const (
	connectTimeoutSeconds = 10
)

// Registry is a registry for nats-leafnodes
type Registry struct {
	credsFilesPath         string
	configFilePath         string
	natsConn               *nats.Conn
	registerSubscription   *nats.Subscription
	unregisterSubscription *nats.Subscription
	state                  *State
	config                 *natsConfig.Config
	natsURI                string
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
		natsURI:        natsURI,
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
		log.Printf("Received register request for network: %s and component: %s\n", creds.Network, creds.Component)
		err = r.addCredentials(creds)
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
		deleteCredsfile, err := r.removeCredentials(creds)
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

func (r *Registry) addCredentials(creds *api.Credentials) error {
	found := false
	for _, remote := range r.config.Leafnodes.Remotes {
		networkName := filepath.Base(remote.Credentials)
		if networkName == fmt.Sprintf("%s.creds", creds.Network) {
			found = true
			break
		}
	}
	if !found {
		err := r.config.AddRemote(creds.NatsAddress, r.credsFile(creds.Network), creds.AccountPublicKey, []string{"local.>"}, []string{"local.>"})
		if err != nil {
			return err
		}
	}
	err := r.state.Update(creds.Network, creds.Component, Add)
	if err != nil {
		return err
	}
	return nil
}

// waitForStreamsDeletion blocks until all the streams are deleted.
func (r *Registry) waitForStreamsDeletion(creds *api.Credentials) error {
	log.Printf("Waiting for streams deletion for network '%s' before leafnode disconnect\n", creds.Network)
	js, err := jetstreams.NewJetstreamControllerWithAddress(creds.Creds, r.natsURI)
	if err != nil {
		return err
	}

	streams, err := js.ListNamesNoDomain()
	if err != nil {
		return err
	}
	log.Printf("Deleting streams for network '%s':\n", creds.Network)
	for _, stream := range streams {
		log.Printf("\t- %s\n", stream)
	}
	err = js.DeleteNoDomain(creds.Network, streams)
	if err != nil {
		return err
	}

	return nil
}

func (r *Registry) removeCredentials(creds *api.Credentials) (bool, error) {
	usage, err := r.state.Usage(creds.Network)
	if err != nil {
		return false, err
	}
	if usage > 0 {
		log.Printf("Removing participant count from network '%s'\n", creds.Network)
		err = r.state.Update(creds.Network, creds.Component, Remove)
		if err != nil {
			return false, err
		}
		usage--
	}
	if usage <= 0 {
		r.waitForStreamsDeletion(creds)
		for _, remote := range r.config.Leafnodes.Remotes {
			if remote.Credentials == r.credsFile(creds.Network) {
				err = r.config.RemoveRemoteByCredsfile(remote.Credentials)
				if err != nil {
					return false, err
				}
				break
			}
		}
		err = r.state.Delete(creds.Network)
		if err != nil {
			return false, err
		}
	}
	return usage <= 0, nil
}

func (r *Registry) credsFile(network string) string {
	return fmt.Sprintf("%s/%s.creds", r.credsFilesPath, network)
}
