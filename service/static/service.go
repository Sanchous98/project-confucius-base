package static

import (
	"github.com/Sanchous98/project-confucius-base/service/web"
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"reflect"
	"sync"
)

type Static struct {
	sync.Mutex
	router *router.Router
	config *config
}

func (s *Static) Make(container src.Container) src.Service {
	s.config = new(config)
	err := s.config.Unmarshall()

	if err != nil {
		panic(err)
	}

	s.router = container.Get(reflect.TypeOf(web.Web{})).(*web.Web).Router
	s.router.ServeFilesCustom("/{filepath:*}", &fasthttp.FS{
		Root:               s.config.Path,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		AcceptByteRange:    true,
		CompressBrotli:     true,
	})

	return s
}
