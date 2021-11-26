package registry

var config = `{
	"pid_file": "/var/run/nats.pid",
	"http": 8222,
	"server_name": "edge",
	"leafnodes": {
		"remotes": []
	},
	"accounts": {
	}
}`

// Registry is a registry for nats-leafnodes
type Registry struct {
	configFileContent string
	credsFilesPath    string
}

// NewRegistry creates a new registry
func NewRegistry() *Registry {
	r := &Registry{
		configFileContent: string(config),
	}
	return r
}
