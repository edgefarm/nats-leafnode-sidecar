package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/edgefarm/nats-leafnode-sidecar/pkg/common"
	"github.com/nats-io/nats.go"

	api "github.com/edgefarm/edgefarm.network/pkg/apis/config/v1alpha1"
)

const (
	connectTimeoutSeconds = 10
)

// Registry is a registry for nats-leafnodes
type Registry struct {
	configFileContent      string
	credsFilesPath         string
	configFilePath         string
	natsConn               *nats.Conn
	registerSubscription   *nats.Subscription
	unregisterSubscription *nats.Subscription
}

// NewRegistry creates a new registry
func NewRegistry(natsConfig string, creds string, natsURI string) (*Registry, error) {
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
	r := &Registry{
		configFileContent: string(config),
		credsFilesPath:    creds,
		configFilePath:    natsConfig,
		natsConn:          nc,
	}
	return r, nil
}

// Start starts the registry and handles all incoming requests for registering and unregistering
func (r *Registry) Start() error {
	var err error
	r.registerSubscription, err = r.natsConn.Subscribe(common.RegisterSubject, func(m *nats.Msg) {
		log.Println("Received register request")
		userCreds := &api.Credentials{}
		err := json.Unmarshal(m.Data, userCreds)
		if err != nil {
			log.Println("Error unmarshalling credentials: ", err)
		}
		err = r.addCredentials(userCreds.UserAccountName, userCreds.Username, userCreds.Password, userCreds.Creds)
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
		err = r.writeFile(fmt.Sprintf("%s/%s.creds", r.credsFilesPath, userCreds.UserAccountName), userCreds.Creds)
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
		userCreds := &api.Credentials{}
		err := json.Unmarshal(m.Data, userCreds)
		if err != nil {
			log.Println("Error unmarshalling credentials: ", err)
		}
		err = r.removeCredentials(userCreds.UserAccountName)
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
		err = r.updateConfigFile()
		if err != nil {
			log.Println(err)
		}
		err = r.removeFile(fmt.Sprintf("%s/%s", r.credsFilesPath, userCreds.UserAccountName))
		if err != nil {
			log.Println(err)
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
