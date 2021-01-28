package static

import (
	"github.com/Sanchous98/project-confucius-base/service/web"
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/valyala/fasthttp"
)

type Static struct {
	web    *web.Web
	config *config
}

func (s *Static) Make(container src.Container) src.Service {
	s.config = new(config)
	err := s.config.Unmarshall()

	if err != nil {
		panic(err)
	}

	s.web = container.Get(&web.Web{}).(*web.Web)
	s.web.Router.ServeFilesCustom("/{filepath:*}", &fasthttp.FS{
		Root:               s.config.Path,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		AcceptByteRange:    true,
		CompressBrotli:     true,
	})

	return s
}
