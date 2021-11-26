package registry

import "os"

// Config returns the current configuration as a JSON string
func (r *Registry) Config() string {
	return r.configFileContent
}

func (r *Registry) updateConfigFile() error {
	file, err := os.Create(r.configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(r.configFileContent)
	if err != nil {
		return err
	}
	return nil
}
