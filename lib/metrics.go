package lib

import (
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	_ "net/http/pprof"
	"sync"
)

type Metrics struct {
	Registry   *prometheus.Registry
	Collectors sync.Map
	web        *Web
}

func (m *Metrics) Construct(web *Web) {
	m.web = web
	m.web.Router.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))
}

func (m *Metrics) Make(container src.Container) src.Service {
	m.web = container.Get(&Web{}).(*Web)
	m.web.Router.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))

	return m
}
