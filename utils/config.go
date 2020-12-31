package utils

import "io/ioutil"

// Config is basic interface for every service configuration
type Config interface {
	HydrateConfig() error
}

// HydrateConfig is a basic function to read configuration from a yaml
func HydrateConfig(c Config, path string, unmarshalMethod func(in []byte, out interface{}) error) (Config, error) {
	config, _ := ioutil.ReadFile(path)
	err := unmarshalMethod(config, c)

	if err != nil {
		return nil, err
	}

	return c, nil
}
