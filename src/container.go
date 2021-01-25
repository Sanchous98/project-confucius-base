package src

import (
	"fmt"
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
		Launch() error
		Shutdown() error
	}

	Launcher interface {
		Launch()
		Shutdown()
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
	s.make()
	err := make(chan error)
	// SIGINT/SIGTERM handling
	osSignals := make(chan os.Signal)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	for _, service := range s.services {
		switch service.service.(type) {
		case Launchable:
			err <- service.service.(Launchable).Launch()
		}
	}

	for {
		select {
		case errors := <-err:
			if errors != nil {
				log.Fatal(errors)
			}
			os.Exit(0)
		case <-osSignals:
			fmt.Printf("\n")
			log.Print("Shutdown signal received.\n")
			s.Shutdown()
			log.Printf("Server gracefully stopped.\n")
		}
	}
}

func (s serviceContainer) Shutdown() {
	err := make(chan error)
	for _, service := range s.services {
		switch service.service.(type) {
		case Launchable:
			err <- service.service.(Launchable).Shutdown()
		}

		select {
		case errors := <-err:
			if errors != nil {
				log.Fatalf("error with graceful close: %s", err)
			}
		}
	}
}

func (s *serviceContainer) make() {
	for _, service := range s.services {
		service.Make(s)
	}
}
