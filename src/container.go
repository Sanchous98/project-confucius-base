package src

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
)

type (
	Service interface {
		sync.Locker
		Make(Container) Service
	}

	Launchable interface {
		Service
		Launch(chan<- error)
		Shutdown(chan<- error)
	}

	Launcher interface {
		Launch()
		Shutdown(chan<- error)
	}

	// Basic Container interface
	Container interface {
		Get(reflect.Type) Service
		Has(reflect.Type) bool
		Set(reflect.Type, Service, bool)
		Launcher
	}

	serviceContainer struct {
		services []*containerEntry
	}
)

func NewContainer() *serviceContainer {
	return &serviceContainer{make([]*containerEntry, 0)}
}

func (s *serviceContainer) Set(abstraction reflect.Type, service Service, singleton bool) {
	s.drop(abstraction)
	s.services = append(s.services, NewEntry(abstraction, service, singleton))
}

func (s *serviceContainer) Get(abstraction reflect.Type) Service {
	for _, service := range s.services {
		if service.Abstraction == abstraction {
			return service.Make(s)
		}
	}

	return nil
}

func (s *serviceContainer) Has(abstraction reflect.Type) bool {
	return s.Get(abstraction) != nil
}

func (s *serviceContainer) drop(abstraction reflect.Type) {
	if !s.Has(abstraction) {
		return
	}

	for index, service := range s.services {
		if service.Abstraction == abstraction {
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
