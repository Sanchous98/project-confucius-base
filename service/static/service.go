package static

import (
	"github.com/Sanchous98/project-confucius-base/service/web"
	"github.com/valyala/fasthttp"
	"log"
	"unsafe"
)

type Static struct {
	web    web.Web
	config *config
}

func (s *Static) Construct(web web.Web) *Static {
	s.config = new(config)
	err := s.config.Unmarshall()
	log.Print(unsafe.Pointer(&web))
	if err != nil {
		panic(err)
	}

	s.web = web
	s.web.Router.ServeFilesCustom("/{filepath:*}", &fasthttp.FS{
		Root:               s.config.Path,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		AcceptByteRange:    true,
		CompressBrotli:     true,
	})

	return s
}
