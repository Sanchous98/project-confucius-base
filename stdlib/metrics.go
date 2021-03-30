package stdlib

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
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
				"/",
				fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()),
			},
		},
		"metrics",
		"/metrics",
	}, "metrics")
}

func (m *Metrics) Destructor() {}
