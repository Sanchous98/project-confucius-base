package src

import (
	"reflect"
)

// Use injectTag to inject dependency into a service
const injectTag = "inject"

var interfaceType = reflect.TypeOf((*Service)(nil)).Elem()

// fillService builds a Service using singletons from Container or new instances of another Services
// TODO: Interface binding
// TODO: Try to eliminate usage of reflect package
func fillService(service Service, container Container) {
	var ok bool
	s := reflect.ValueOf(service).Elem()

	for i := 0; i < s.NumField(); i++ {
		var newService reflect.Value
		field := s.Field(i)
		dependencyType := field.Type()

		if container.Has(dependencyType) {
			newService = reflect.ValueOf(container.Get(dependencyType))
		} else if dependencyType.Implements(interfaceType) {
			newService = reflect.New(dependencyType.Elem())
			fillService(newService.Interface().(Service), container)
		} else {
			continue
		}

		newService.Interface().(Service).Constructor()
		_, ok = s.Type().Field(i).Tag.Lookup(injectTag)

		if ok {
			field.Set(newService)
		}
	}

	service.Constructor()
}
