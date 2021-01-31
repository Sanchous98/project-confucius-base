package lib

import (
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/Sanchous98/project-confucius-base/utils"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
)

const staticConfigPath = "config/static.yaml"

type (
	staticConfig struct {
		Path string
	}

	Static struct {
		web    *Web
		config *staticConfig
	}
)

func (c *staticConfig) Unmarshall() error {
	absPath, _ := filepath.Abs(staticConfigPath)
	content, err := ioutil.ReadFile(absPath)
	cfg, err := utils.HydrateConfig(c, content, yaml.Unmarshal)

	if err != nil {
		return err
	}

	c = cfg.(*staticConfig)

	return nil
}

func (s *Static) Make(container src.Container) src.Service {
	s.config = new(staticConfig)
	err := s.config.Unmarshall()

	if err != nil {
		panic(err)
	}

	s.web = container.Get(&Web{}).(*Web)
	s.web.Router.ServeFilesCustom("/{filepath:*}", &fasthttp.FS{
		Root:               s.config.Path,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		AcceptByteRange:    true,
		CompressBrotli:     true,
	})

	return s
}
