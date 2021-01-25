package graphql

import (
	"github.com/Sanchous98/project-confucius-base/utils"
	"gopkg.in/yaml.v3"
)

const ConfigPath = "config/graphql.yaml"

type config struct {
	SchemaPath string `yaml:"schema_path"`
}

func (gqlc *config) Unmarshall() error {
	cfg, err := utils.HydrateConfig(gqlc, ConfigPath, yaml.Unmarshal)

	if err != nil {
		return err
	}

	gqlc = cfg.(*config)

	return nil
}
