package stdlib

import (
	"github.com/Sanchous98/project-confucius-base/utils"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

const staticConfigPath = "config/static.yaml"

type (
	staticConfig struct {
		Path     string
		Compress bool
		Indexes  []string
	}

	Static struct {
		Web    *Web `inject:""`
		Log    *Log `inject:""`
		config *staticConfig
	}
)

func (c *staticConfig) Unmarshall() error {
	return utils.Unmarshall(c, staticConfigPath, yaml.Unmarshal)
}

func (s *Static) Constructor() {
	s.config = new(staticConfig)
	err := s.config.Unmarshall()

	if err != nil {
		s.Log.Critical(err)
	}

	path, err := filepath.Abs(s.config.Path)

	fs := &fasthttp.FS{
		Root:               path,
		IndexNames:         s.config.Indexes,
		GenerateIndexPages: false,
		AcceptByteRange:    true,
		CompressBrotli:     s.config.Compress,
	}

	s.Web.AddEntryPoint(&EntryPoint{
		[]*Route{{
			MethodGet,
			"/{filepath:*}",
			fs.NewRequestHandler(),
		}},
		"", "",
	}, "")
}
