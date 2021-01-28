package src

import (
	"reflect"
	"sync"
)

// containerEntry is a wrapper for services to use binding by abstraction and singletons
type containerEntry struct {
	sync.RWMutex
	sync.Once
	service     Service
	Abstraction reflect.Type
}

func NewEntry(abstraction reflect.Type, service Service) *containerEntry {
	return &containerEntry{
		Abstraction: abstraction,
		service:     service,
	}
}

func (c *containerEntry) Make(container Container) Service {
	c.Do(func() {
		c.RLock()
		defer c.RUnlock()

		c.service.Make(container)
	})

	return c.service
}
