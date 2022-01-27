package registry

import (
	"os"
	"path/filepath"
)

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

func (r *Registry) writeFile(path string, content string) error {
	newpath := filepath.Join(filepath.Dir(path))
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func (r *Registry) removeFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content := ""
	_, err = file.Read([]byte(content))
	if err != nil {
		return "", err
	}
	return content, nil
}
