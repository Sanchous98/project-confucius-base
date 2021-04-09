package src

import (
	"reflect"
)

// Use injectTag to inject dependency into a service
const injectTag = "inject"

// If dependency is not found in container, we make a new instance for every struct field with "inject" tag
// cache is used to avoid calling reflect multiple types, because we can just copy once instantiated services with a
// default configuration
var cache = make(map[string]Service)

// fillService builds a Service using singletons from Container or new instances of another Services
// TODO: Try to eliminate usage of reflect package
func fillService(service Service, container Container) {
	ok := false
	s := reflect.ValueOf(service).Elem()

	for i := 0; i < s.NumField(); i++ {
		_, ok = s.Type().Field(i).Tag.Lookup(injectTag)

		if !ok {
			continue
		}

		var newService Service
		field := s.Field(i)
		dependencyType := field.Type()

		if container != nil && container.Has(dependencyType) {
			// If service is bound, take it from the container
			newService = container.Get(dependencyType)
		} else if _, ok = field.Interface().(Service); ok {
			// If service is not bound, find a cached instance
			if serv, ok := cache[dependencyType.String()]; ok {
				newService = *copyService(&serv)
			} else {
				// If service is not found even in cache, inject a new instance
				newService = reflect.New(dependencyType.Elem()).Interface().(Service)
				cache[dependencyType.String()] = newService
				fillService(newService, container)
			}
		} else {
			// If service is not found, just skip it
			continue
		}

		newService.Constructor()
		field.Set(reflect.ValueOf(newService))
	}
	service.Constructor()
}

func copyService(service *Service) *Service {
	serv := *service
	return &serv
}
