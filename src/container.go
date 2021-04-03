package src

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

type (
	// Service represents an entity in a microservice architecture
	Service interface {
		Constructor()
		Destructor()
	}

	// Launchable Base interface. Represents a Service that can be launched in a separate thread
	Launchable interface {
		// Launch method is a service coroutine
		Launch(chan<- error)
		// Shutdown terminates service coroutine
		Shutdown(chan<- error)
	}

	// Launcher Base interface. Represent entity that launches and shuts down services
	Launcher interface {
		Launch()
		Shutdown(chan<- error)
	}

	// Container Base interface. Conducts services
	Container interface {
		Get(reflect.Type) Service
		Has(reflect.Type) bool
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

// Set bins a singleton service.
// If you want to use a new instance on every binding, container will resolve it automatically
func (s *serviceContainer) Set(service Service) {
	s.services = append(s.services, NewEntry(service))
}

// Get returns bound service
func (s *serviceContainer) Get(abstraction reflect.Type) Service {
	for _, service := range s.services {
		if reflect.TypeOf(service.service) == abstraction {
			return service.Make(s)
		}
	}

	return nil
}

// Has checks, if a service is bound
func (s *serviceContainer) Has(abstraction reflect.Type) bool {
	return s.Get(abstraction) != nil
}

func (s *serviceContainer) drop(abstraction reflect.Type) {
	if !s.Has(abstraction) {
		return
	}

	for index, service := range s.services {
		if reflect.TypeOf(service.service) == abstraction {
			s.services = append(s.services[:index], s.services[index+1:]...)
			return
		}
	}
}

// TODO: Move into app
func (s *serviceContainer) Launch() {
	// TODO: Sort dependencies
	err := make(chan error)
	osSignals := make(chan os.Signal)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
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
