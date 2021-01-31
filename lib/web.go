package lib

import (
	"crypto/tls"
	"fmt"
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/Sanchous98/project-confucius-base/utils"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net"
	"path/filepath"
)

const webConfigPath = "config/web.yaml"

type (
	webConfig struct {
		CertsPath   string `yaml:"certs_path"`
		Addr        string
		Port        uint8
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
	}
)

func (c *webConfig) Unmarshall() error {
	absPath, _ := filepath.Abs(webConfigPath)
	content, err := ioutil.ReadFile(absPath)
	cfg, err := utils.HydrateConfig(c, content, yaml.Unmarshal)

	if err != nil {
		return err
	}

	c = cfg.(*webConfig)

	return err
}

func (c *webConfig) getFullAddress() string {
	return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}

func (w *Web) Make(src.Container) src.Service {
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

	return w
}

func (w *Web) Launch(err chan<- error) {
	// Let's Encrypt tls-alpn-01 only works on Port 443.
	ln, e := net.Listen("tcp4", w.config.getFullAddress())

	if e != nil {
		err <- e
	}

	lnTls := tls.NewListener(ln, w.TLSConfig)
	if w.config.Compression.Enabled {
		fasthttp.CompressHandlerBrotliLevel(w.Router.Handler, w.config.Compression.Level, w.config.Compression.Level)
	}

	w.Server = &fasthttp.Server{
		Handler: w.Router.Handler,
	}

	log.Print("Server started")
	err <- w.Server.Serve(lnTls)
}

func (w *Web) Shutdown(chan<- error) {
	w.Server.DisableKeepalive = true
}
