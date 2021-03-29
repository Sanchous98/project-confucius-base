package utils

import (
	"io/ioutil"
	"path/filepath"
)

type unmarshallMethod func(in []byte, out interface{}) error

type (
	// ConfigType is basic interface for every service configuration
	ConfigType interface {
		Unmarshall() error
	}
)

// Unmarshall is a basic function to read configuration from a yaml
func Unmarshall(config ConfigType, path string, unmarshallMethod unmarshallMethod) error {
	absPath, _ := filepath.Abs(path)
	content, err := ioutil.ReadFile(absPath)

	if err != nil {
		return err
	}

	err = unmarshallMethod(content, config)

	return err
}
