package src

import (
	"reflect"
)

// Use injectTag to inject dependency into a service
const injectTag = "inject"

// interfaceType is Service interface null pointer
var interfaceType = reflect.TypeOf((*Service)(nil)).Elem()

// fillService builds a Service using singletons from Container or new instances of another Services
// TODO: Try to eliminate usage of reflect package
func fillService(service Service, container Container) {
	ok := false
	s := reflect.ValueOf(service).Elem()

	for i := 0; i < s.NumField(); i++ {
		var newService reflect.Value
		field := s.Field(i)
		dependencyType := field.Type()
		_, ok = s.Type().Field(i).Tag.Lookup(injectTag)

		if !ok {
			continue
		}

		if container != nil && container.Has(dependencyType) {
			// If service is bound, take it from the container
			newService = reflect.ValueOf(container.Get(dependencyType))
		} else if dependencyType.Implements(interfaceType) {
			// If service is not bound, inject a new instance
			newService = reflect.New(dependencyType.Elem())
			fillService(newService.Interface().(Service), container)
		} else {
			// If service is not found, just skip it
			continue
		}

		newService.Interface().(Service).Constructor()
		field.Set(newService)
	}

	service.Constructor()
}
