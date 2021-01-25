package src

import (
	"reflect"
	"sync"
)

// containerEntry is a wrapper for services to use binding by abstraction and singletons
type containerEntry struct {
	sync.Mutex
	service             Service
	IsSingleton, isMade bool
	Abstraction         reflect.Type
}

func NewEntry(abstraction reflect.Type, service Service, isSingleton bool) *containerEntry {
	return &containerEntry{
		Abstraction: abstraction,
		service:     service,
		IsSingleton: isSingleton,
	}
}

func (c *containerEntry) Make(container Container) Service {
	if !c.IsSingleton {
		c.service = reflect.New(reflect.ValueOf(c.service).Elem().Type()).Interface().(Service)

		return c.service.Make(container)
	}

	if !c.isMade {
		c.Lock()
		c.service.Make(container)
		c.isMade = true
		c.Unlock()
	}

	return c.service
}
