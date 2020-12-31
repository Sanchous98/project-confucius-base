package confucius

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Sanchous98/project-confucius-base/utils"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type (
	// Basic Service interface
	Service interface {
		Serve(handler *mux.Router) error
		Stop()
		Init() error
	}

	// Basic Container interface
	Container interface {
		Get(service string) (*Service, Status)
		Set(name string, service Service)
		Has(service string) bool
		Service
	}

	serviceContainer struct {
		sync.Mutex
		config   *Config
		server   *http.Server
		services []*containerEntry
	}
)

func NewContainer(config utils.Config) Container {
	return &serviceContainer{
		services: make([]*containerEntry, 0),
		config:   config.(*Config),
		server:   &http.Server{Addr: ":80"},
	}
}

func (s *serviceContainer) Get(service string) (*Service, Status) {
	s.Lock()
	defer s.Unlock()

	for _, e := range s.services {
		if e.name == service {
			return &e.service, e.getStatus()
		}
	}

	return nil, Undefined
}

func (s *serviceContainer) Set(name string, service Service) {
	s.Lock()
	s.services = append(s.services, &containerEntry{
		name:    name,
		service: service,
		status:  Inactive,
	})
	s.Unlock()
}

func (s *serviceContainer) Has(service string) bool {
	s.Lock()
	defer s.Unlock()

	for _, entry := range s.services {
		if entry.name == service {
			return true
		}
	}

	return false
}

// Serve launches all services
func (s *serviceContainer) Serve(router *mux.Router) error {
	for _, entry := range s.services {
		if entry.hasStatus(Ok) {
			entry.setStatus(Serving)
			if err := entry.service.Serve(router); err != nil {
				return errors.Wrap(err, fmt.Sprintf("[%s]", entry.name))
			}
		}
	}

	return s.selfServe(router)
}

// Stop shuts down all services
func (s *serviceContainer) Stop() {
	for _, entry := range s.services {
		if entry.hasStatus(Serving) {
			entry.setStatus(Stopping)
			entry.service.Stop()
			entry.setStatus(Stopped)
		}
	}

	s.selfStop()
}

// Init initializes all services
func (s *serviceContainer) Init() error {
	for _, entry := range s.services {
		if err := initService(entry); err != nil {
			return err
		}
	}

	return s.selfInit()
}

func initService(entry *containerEntry) error {
	if entry.getStatus() >= Ok {
		return fmt.Errorf("service [%s] has already been configured", entry.name)
	}

	err := entry.service.Init()

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("service [%s] cannot be initialized", entry.name))
	}

	entry.setStatus(Ok)

	return nil
}

func (s *serviceContainer) selfInit() error {
	return s.config.HydrateConfig()
}

func (s *serviceContainer) selfServe(router *mux.Router) error {
	// TODO: Test on real machine, because not working for localhost
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(os.Getenv("CERTS_PATH")),
		HostPolicy: autocert.HostWhitelist(s.config.Domain),
	}

	s.server.Handler = certManager.HTTPHandler(router)
	s.server.TLSConfig = &tls.Config{
		GetCertificate: certManager.GetCertificate,
	}

	log.Print("SERVER STARTED")
	log.Fatal(s.server.ListenAndServe())

	return nil
}

func (s *serviceContainer) selfStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_ = s.server.Shutdown(ctx)
	cancel()
}
