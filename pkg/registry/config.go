package registry

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func (r *Registry) updateConfigFile() error {
	file, err := os.Create(r.configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	j, err := r.config.ToJSON()
	if err != nil {
		return err
	}
	_, err = file.WriteString(j)
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
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
