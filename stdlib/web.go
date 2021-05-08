package stdlib

import (
	"crypto/tls"
	"fmt"
	"github.com/Sanchous98/project-confucius-base/utils"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"net/http"
	"sync"
)

const webConfigPath = "config/web.yaml"

const (
	MethodGet Method = iota
	MethodHead
	MethodPost
	MethodPut
	MethodPatch
	MethodDelete
	MethodConnect
	MethodOptions
	MethodTrace
)

type (
	// Method type is an integer version of http package methods, because switch works faster on integers
	Method uint8

	// Middleware makes some actions before request handling
	Middleware func(fasthttp.RequestHandler) fasthttp.RequestHandler

	webConfig struct {
		CertsPath   string `yaml:"certs_path"`
		Addr        string
		Port        uint16
		Compression struct {
			Enabled bool
			Level   int
		}
		Whitelist []string `yaml:"whitelist"`
	}

	// Web is a main http server
	Web struct {
		config      *webConfig
		router      *router.Router
		server      *fasthttp.Server
		certManager *autocert.Manager
		tlsConfig   *tls.Config
		Log         *Log `inject:""`
		entryPoints map[string][]*EntryPoint
		sync.Once
	}

	// Route is an abstraction over fasthttp routes. Contains HTTP Method, path and handler
	Route struct {
		Method  Method
		Path    string
		Handler fasthttp.RequestHandler
	}

	// EntryPoint is an abstraction over Route. Represents a collection of routes, united by group and prefix
	EntryPoint struct {
		Routes      []*Route
		Name        string
		RoutePrefix string
		Middlewares []Middleware
	}
)

func NewEntryPoint(name, prefix string) *EntryPoint {
	return &EntryPoint{make([]*Route, 0), name, prefix, make([]Middleware, 0)}
}

func (e *EntryPoint) AddRoute(route *Route) {
	e.Routes = append(e.Routes, route)
}

func (e *EntryPoint) AddMiddleware(middleware Middleware) {
	e.Middlewares = append(e.Middlewares, middleware)
}

func (m Method) String() string {
	switch m {
	case MethodGet:
		return http.MethodGet
	case MethodHead:
		return http.MethodHead
	case MethodPost:
		return http.MethodPost
	case MethodPut:
		return http.MethodPut
	case MethodPatch:
		return http.MethodPatch
	case MethodDelete:
		return http.MethodDelete
	case MethodConnect:
		return http.MethodConnect
	case MethodOptions:
		return http.MethodOptions
	case MethodTrace:
		return http.MethodTrace
	default:
		return ""
	}
}

func (c *webConfig) Unmarshall() error {
	return utils.Unmarshall(c, webConfigPath, yaml.Unmarshal)
}

func (c *webConfig) getFullAddress() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}

// Launch web server
func (w *Web) Launch() {
	w.server.Handler = w.router.Handler

	w.Do(func() {
		for _, entryPoints := range w.entryPoints {
			for _, entryPoint := range entryPoints {
				for _, route := range entryPoint.Routes {
					handler := route.Handler

					for _, middleware := range entryPoint.Middlewares {
						handler = middleware(handler)
					}

					if route.Path[0] != byte('/') {
						w.Log.Error(
							fmt.Errorf(
								"path must begin with '/' in entry point '%s', path '%s'. Skipping route",
								entryPoint.Name,
								route.Path,
							),
						)
						continue
					}
					switch route.Method {
					case MethodGet:
						w.router.GET(entryPoint.RoutePrefix+route.Path, handler)
					case MethodHead:
						w.router.HEAD(entryPoint.RoutePrefix+route.Path, handler)
					case MethodPost:
						w.router.POST(entryPoint.RoutePrefix+route.Path, handler)
					case MethodPut:
						w.router.PUT(entryPoint.RoutePrefix+route.Path, handler)
					case MethodPatch:
						w.router.PATCH(entryPoint.RoutePrefix+route.Path, handler)
					case MethodDelete:
						w.router.DELETE(entryPoint.RoutePrefix+route.Path, handler)
					case MethodConnect:
						w.router.CONNECT(entryPoint.RoutePrefix+route.Path, handler)
					case MethodOptions:
						w.router.OPTIONS(entryPoint.RoutePrefix+route.Path, handler)
					case MethodTrace:
						w.router.TRACE(entryPoint.RoutePrefix+route.Path, handler)
					}
				}
			}
		}
	})

	if w.config.Compression.Enabled {
		fasthttp.CompressHandlerBrotliLevel(w.router.Handler, w.config.Compression.Level, w.config.Compression.Level)
	}

	// Let's Encrypt tls-alpn-01 only works on Port 443.
	if w.config.Port == 443 {
		ln, e := net.Listen("tcp4", w.config.getFullAddress())

		if e != nil {
			panic(e)
		}

		log.Printf("Server started on https://%s", w.config.getFullAddress())
		lnTls := tls.NewListener(ln, w.tlsConfig)

		if e = w.server.Serve(lnTls); e != nil {
			w.Log.Alert(e)
		}
	} else {
		log.Printf("Server started on http://%s", w.config.getFullAddress())
		if e := w.server.ListenAndServe(w.config.getFullAddress()); e != nil {
			w.Log.Alert(e)
		}
	}
}

// Shutdown web server
func (w *Web) Shutdown() {
	w.server.DisableKeepalive = true
	w.Log.Info(w.server.Shutdown())
}

func (w *Web) Constructor() {
	w.config = new(webConfig)
	err := w.config.Unmarshall()
	if err != nil {
		w.Log.Emergency(err)
	}
	w.router = router.New()
	w.server = &fasthttp.Server{}

	w.certManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(w.config.Whitelist...),
		Cache:      autocert.DirCache(w.config.CertsPath),
	}

	w.tlsConfig = &tls.Config{
		GetCertificate: w.certManager.GetCertificate,
		NextProtos:     []string{"http/1.1", acme.ALPNProto},
	}
}

func (w *Web) AddEntryPoint(entryPoint *EntryPoint, group string) {
	if w.entryPoints == nil {
		w.entryPoints = make(map[string][]*EntryPoint)
	}

	w.entryPoints[group] = append(w.entryPoints[group], entryPoint)
}

func (w *Web) DropEntryPoint(name string, group string) {
	for index, entryPoint := range w.entryPoints[group] {
		if entryPoint.Name == name {
			w.entryPoints[group] = append(w.entryPoints[group][:index], w.entryPoints[group][index+1:]...)
			return
		}
	}
}

func (w *Web) EntryPointExists(name string, group string) bool {
	for _, entryPoint := range w.entryPoints[group] {
		if entryPoint.Name == name {
			return true
		}
	}

	return false
}
