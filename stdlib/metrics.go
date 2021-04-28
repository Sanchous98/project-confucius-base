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
	entrypoint := NewEntryPoint("debug", "/debug")
	entrypoint.AddRoute(&Route{
		MethodGet,
		"/prometheus",
		fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()),
	})
	entrypoint.AddRoute(&Route{
		MethodGet,
		"/pprof/{path:*}",
		pprofhandler.PprofHandler,
	})

	m.Web.AddEntryPoint(entrypoint, "debug")
}
