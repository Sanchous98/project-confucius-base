package web

import (
	"fmt"
	"github.com/Sanchous98/project-confucius-base/utils"
	"gopkg.in/yaml.v3"
)

const ConfigPath = "config/web.yaml"

type config struct {
	CertsPath   string `yaml:"certs_path"`
	Addr        string
	Port        uint8
	Compression struct {
		Enabled bool
		Level   int
	}
	Whitelist []string `yaml:"whitelist"`
}

func (c *config) Unmarshall() error {
	cfg, err := utils.HydrateConfig(c, ConfigPath, yaml.Unmarshal)

	if err != nil {
		return err
	}

	c = cfg.(*config)

	return err
}

func (c *config) getFullAddress() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}
