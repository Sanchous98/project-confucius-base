package static

import (
	"github.com/Sanchous98/project-confucius-base/utils"
	"gopkg.in/yaml.v3"
)

const ConfigPath = "config/static.yaml"

type config struct {
	Path string
}

func (c *config) Unmarshall() error {
	cfg, err := utils.HydrateConfig(c, ConfigPath, yaml.Unmarshal)

	if err != nil {
		return err
	}

	c = cfg.(*config)

	return nil
}
