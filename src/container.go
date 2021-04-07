package src

import (
	"reflect"
)

type (
	// Service represents an entity in a microservice architecture
	Service interface {
		Constructor()
	}

	// Launcher Base interface. Represent entity that launches and shuts down services
	Launcher interface {
		Launch()
		Shutdown()
	}

	// Container Base interface. Conducts services
	Container interface {
		Get(reflect.Type) Service
		Has(reflect.Type) bool
		Set(Service)
		GetServices() []Service
	}

	// TODO: Make serviceContainer implement binary tree data structure
	serviceContainer struct {
		services []*containerEntry
	}
)

func NewContainer() *serviceContainer {
	return &serviceContainer{make([]*containerEntry, 0)}
}

// Set binds a singleton service.
// If you want to use a new instance on every binding, container will resolve it automatically
func (s *serviceContainer) Set(service Service) {
	s.services = append(s.services, NewEntry(service))
}

// Get returns bound service
func (s *serviceContainer) Get(abstraction reflect.Type) Service {
	for _, service := range s.services {
		if reflect.TypeOf(service.service) == abstraction {
			return service.make(s)
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

// TODO: Sort dependencies
func (s *serviceContainer) GetServices() []Service {
	mappedServices := make([]Service, len(s.services))

	for index, service := range s.services {
		mappedServices[index] = service.make(s)
	}

	return mappedServices
}
