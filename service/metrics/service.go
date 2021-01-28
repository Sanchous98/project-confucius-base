package metrics

import (
	"github.com/Sanchous98/project-confucius-base/service/web"
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"sync"
)

type Metrics struct {
	sync.Mutex
	Registry   *prometheus.Registry
	Collectors sync.Map
	web        *web.Web
}

func (m *Metrics) Construct(container src.Container) *Metrics {
	container.Inject(m.web)
	m.web.Router.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))

	return m
}
