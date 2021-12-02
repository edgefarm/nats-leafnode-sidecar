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
	configFilePath        = "/config/nats.json"
	connectTimeoutSeconds = 10
)

// Registry is a registry for nats-leafnodes
type Registry struct {
	configFileContent string
	credsFilesPath    string
	configFilePath    string
	natsConn          *nats.Conn
}

// NewRegistry creates a new registry
func NewRegistry(creds string, natsURI string) (*Registry, error) {
	opts := []nats.Option{nats.Timeout(time.Duration(1) * time.Second)}
	opts = setupConnOptions(opts)
	ncChan := make(chan *nats.Conn)
	go func() {
		for {
			fmt.Printf("\rConnecting to nats server: %s\n", natsURI)
			nc, err := nats.Connect(natsURI, opts...)
			if err != nil {
				fmt.Printf("Connect failed to %s: %s\n", natsURI, err)
			} else {
				fmt.Printf("Connected to '%s'\n", natsURI)
				ncChan <- nc
				return
			}
			func() {
				for i := connectTimeoutSeconds; i >= 0; i-- {
					time.Sleep(time.Second)
					fmt.Printf("\rReconnecting in %2d seconds", i)
				}
				fmt.Println("")
			}()
		}
	}()

	nc := <-ncChan
	r := &Registry{
		configFileContent: string(config),
		credsFilesPath:    creds,
		configFilePath:    configFilePath,
		natsConn:          nc,
	}
	return r, nil
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := 2 * time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
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
