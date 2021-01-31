package utils

// Config is basic interface for every service configuration
type Config interface {
	Unmarshall() error
}

// HydrateConfig is a basic function to read configuration from a yaml
func HydrateConfig(c Config, content []byte, unmarshalMethod func(in []byte, out interface{}) error) (Config, error) {
	err := unmarshalMethod(content, c)

	if err != nil {
		return nil, err
	}

	return c, nil
}
