package web

import (
	"crypto/tls"
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net"
	"sync"
)

type Web struct {
	sync.Mutex
	config      *config
	Server      *fasthttp.Server
	Router      *router.Router
	CertManager *autocert.Manager
	TLSConfig   *tls.Config
}

func (w *Web) Make(src.Container) src.Service {
	w.config = new(config)
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

func (w *Web) Launch() error {
	// Let's Encrypt tls-alpn-01 only works on Port 443.
	ln, _ := net.Listen("tcp4", w.config.getFullAddress())

	lnTls := tls.NewListener(ln, w.TLSConfig)
	if w.config.Compression.Enabled {
		fasthttp.CompressHandlerBrotliLevel(w.Router.Handler, w.config.Compression.Level, w.config.Compression.Level)
	}

	w.Server = &fasthttp.Server{
		Handler: w.Router.Handler,
	}

	log.Print("Server started")
	return w.Server.Serve(lnTls)
}

func (w *Web) Shutdown() error {
	w.Server.DisableKeepalive = true

	return nil
}
