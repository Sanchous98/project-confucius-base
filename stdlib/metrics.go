package stdlib

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/valyala/fasthttp/pprofhandler"
	_ "net/http/pprof"
	"sync"
)

type Metrics struct {
	Registry   *prometheus.Registry
	Collectors sync.Map
	Web        *Web `inject:""`
}

func (m *Metrics) Constructor() {
	m.Web.AddEntryPoint(&EntryPoint{
		[]*Route{
			{
				MethodGet,
				"/prometheus",
				fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()),
			},
			{
				MethodGet,
				"/pprof/{path:*}",
				pprofhandler.PprofHandler,
			},
		},
		"debug",
		"/debug",
	}, "debug")
}
