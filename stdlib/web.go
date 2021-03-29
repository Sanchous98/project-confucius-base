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
)

const webConfigPath = "config/web.yaml"

type (
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
		Server      *fasthttp.Server
		Router      *router.Router
		CertManager *autocert.Manager
		TLSConfig   *tls.Config
		Log         *Log `inject:""`
	}
)

func (c *webConfig) Unmarshall() error {
	return utils.Unmarshall(c, webConfigPath, yaml.Unmarshal)
}

func (c *webConfig) getFullAddress() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}

func (w *Web) Launch(err chan<- error) {
	if w.config.Compression.Enabled {
		fasthttp.CompressHandlerBrotliLevel(w.Router.Handler, w.config.Compression.Level, w.config.Compression.Level)
	}

	w.Server = &fasthttp.Server{
		Handler: w.Router.Handler,
	}

	log.Print("Server started")

	// Let's Encrypt tls-alpn-01 only works on Port 443.
	if w.config.Port == 443 {
		ln, e := net.Listen("tcp4", w.config.getFullAddress())

		if e != nil {
			err <- e
		}

		lnTls := tls.NewListener(ln, w.TLSConfig)
		err <- w.Server.Serve(lnTls)
	} else {
		err <- w.Server.ListenAndServe(w.config.getFullAddress())
	}
}

func (w *Web) Shutdown(chan<- error) {
	w.Server.DisableKeepalive = true
}

func (w *Web) Constructor() {
	w.config = new(webConfig)
	_ = w.config.Unmarshall()

	w.Router = router.New()
	w.CertManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(w.config.Whitelist...),
		Cache:      autocert.DirCache(w.config.CertsPath),
	}

	w.TLSConfig = &tls.Config{
		GetCertificate: w.CertManager.GetCertificate,
		NextProtos:     []string{"http/1.1", acme.ALPNProto},
	}
}

func (w *Web) Destructor() {}
