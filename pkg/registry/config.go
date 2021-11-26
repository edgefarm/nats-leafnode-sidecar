package registry

const (
	credsFilesPath = "/creds"
	configFilePath = "/nats.json"
)

// Config returns the current configuration as a JSON string
func (r *Registry) Config() string {
	return r.configFileContent
}
