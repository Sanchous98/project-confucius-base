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
	Method uint8

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

	Web struct {
		config      *webConfig
		router      *router.Router
		server      *fasthttp.Server
		certManager *autocert.Manager
		tlsConfig   *tls.Config
		Log         *Log `inject:""`
		entryPoints map[string][]*EntryPoint
	}

	Route struct {
		Method  Method
		Path    string
		Handler fasthttp.RequestHandler
	}

	EntryPoint struct {
		Routes      []*Route
		Name        string
		RoutePrefix string
	}
)

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

func (w *Web) Launch(err chan<- error) {
	for group, entryPoints := range w.entryPoints {
		for _, entryPoint := range entryPoints {
			for _, route := range entryPoint.Routes {
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
					w.router.GET("/"+group+route.Path, route.Handler)
				case MethodHead:
					w.router.HEAD("/"+group+route.Path, route.Handler)
				case MethodPost:
					w.router.POST("/"+group+route.Path, route.Handler)
				case MethodPut:
					w.router.PUT("/"+group+route.Path, route.Handler)
				case MethodPatch:
					w.router.PATCH("/"+group+route.Path, route.Handler)
				case MethodDelete:
					w.router.DELETE("/"+group+route.Path, route.Handler)
				case MethodConnect:
					w.router.CONNECT("/"+group+route.Path, route.Handler)
				case MethodOptions:
					w.router.OPTIONS("/"+group+route.Path, route.Handler)
				case MethodTrace:
					w.router.TRACE("/"+group+route.Path, route.Handler)
				}
			}
		}
	}

	if w.config.Compression.Enabled {
		fasthttp.CompressHandlerBrotliLevel(w.router.Handler, w.config.Compression.Level, w.config.Compression.Level)
	}

	w.server = &fasthttp.Server{
		Handler: w.router.Handler,
	}

	log.Print("Server started")

	// Let's Encrypt tls-alpn-01 only works on Port 443.
	if w.config.Port == 443 {
		ln, e := net.Listen("tcp4", w.config.getFullAddress())

		if e != nil {
			err <- e
		}

		lnTls := tls.NewListener(ln, w.tlsConfig)
		err <- w.server.Serve(lnTls)
	} else {
		err <- w.server.ListenAndServe(w.config.getFullAddress())
	}
}

func (w *Web) Shutdown(chan<- error) {
	w.server.DisableKeepalive = true
}

func (w *Web) Constructor() {
	w.config = new(webConfig)
	_ = w.config.Unmarshall()

	w.router = router.New()
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

func (w *Web) Destructor() {}

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
