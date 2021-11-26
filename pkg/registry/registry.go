package registry

import (
	"fmt"
	"os"
)

const (
	credsFilesPath = "/creds"
	configFilePath = "/nats.json"
)

// Registry is a registry for nats-leafnodes
type Registry struct {
	configFileContent string
	credsFilesPath    string
	configFilePath    string
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	r := &Registry{
		configFileContent: string(config),
		credsFilesPath:    credsFilesPath,
		configFilePath:    configFilePath,
	}
	return r
}

// Shutdown shuts down the registry
func (r *Registry) Shutdown() {
	fmt.Println("Shutting down registry")
	os.Exit(0)
}