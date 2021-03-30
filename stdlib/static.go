package stdlib

import (
	"github.com/Sanchous98/project-confucius-base/utils"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
)

const staticConfigPath = "config/static.yaml"

type (
	staticConfig struct {
		Path string
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

	fs := &fasthttp.FS{
		Root:               s.config.Path,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		AcceptByteRange:    true,
		CompressBrotli:     true,
	}

	fs.PathRewrite = fasthttp.NewVHostPathRewriter(0)
	s.Web.router.ServeFilesCustom("/{filepath:*}", fs)
}

func (s *Static) Destructor() {}
