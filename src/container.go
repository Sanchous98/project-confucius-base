package src

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

type (
	Service interface {
		Make(Container) Service
	}

	Launchable interface {
		Launch(chan<- error)
		Shutdown(chan<- error)
	}

	Launcher interface {
		Launch()
		Shutdown(chan<- error)
	}

	// Basic Container interface
	Container interface {
		Get(Service) Service
		Has(Service) bool
		Set(Service)
		Launcher
	}

	serviceContainer struct {
		services []*containerEntry
	}
)

func NewContainer() *serviceContainer {
	return &serviceContainer{make([]*containerEntry, 0)}
}

func (s *serviceContainer) Set(service Service) {
	s.services = append(s.services, NewEntry(reflect.TypeOf(service), service))
}

func (s *serviceContainer) Get(abstraction Service) Service {
	for _, service := range s.services {
		if service.Abstraction == reflect.TypeOf(abstraction) {
			return service.Make(s)
		}
	}

	return nil
}

func (s *serviceContainer) Has(abstraction Service) bool {
	return s.Get(abstraction) != nil
}

func (s *serviceContainer) drop(abstraction Service) {
	if !s.Has(abstraction) {
		return
	}

	for index, service := range s.services {
		if service.Abstraction == reflect.TypeOf(abstraction) {
			s.services = append(s.services[:index], s.services[index+1:]...)
			return
		}
	}
}

func (s *serviceContainer) Launch() {
	err := make(chan error)
	osSignals := make(chan os.Signal)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
	for _, service := range s.services {
		service.Make(s)
	}

	for _, service := range s.services {
		switch service.service.(type) {
		case Launchable:
			go service.service.(Launchable).Launch(err)
		}
	}

	select {
	case <-osSignals:
		log.Print("\nShutdown signal received.")
		s.Shutdown(err)
		log.Print("\nServer gracefully stopped.\n")
		os.Exit(0)
	case errors := <-err:
		if errors != nil {
			s.Shutdown(err)
			log.Fatal(errors)
		}
	}
}

func (s serviceContainer) Shutdown(err chan<- error) {
	for _, service := range s.services {
		switch service.service.(type) {
		case Launchable:
			go service.service.(Launchable).Shutdown(err)
		}
	}
}
