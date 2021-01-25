package utils

import (
	"io/ioutil"
	"path/filepath"
)

// Config is basic interface for every service configuration
type Config interface {
	Unmarshall() error
}

// HydrateConfig is a basic function to read configuration from a yaml
func HydrateConfig(c Config, path string, unmarshalMethod func(in []byte, out interface{}) error) (Config, error) {
	absPath, _ := filepath.Abs(path)
	config, err := ioutil.ReadFile(absPath)

	if err != nil {
		return nil, err
	}
	err = unmarshalMethod(config, c)

	if err != nil {
		return nil, err
	}

	return c, nil
}
