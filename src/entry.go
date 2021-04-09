package src

import "sync"

// containerEntry is a wrapper for bound services. Makes a singleton
type containerEntry struct {
	sync.RWMutex
	sync.Once
	service Service
}

func NewEntry(service Service) *containerEntry {
	return &containerEntry{service: service}
}

func (c *containerEntry) make(container Container) Service {
	c.Do(func() {
		c.RLock()
		fillService(c.service, container)
		c.RUnlock()
	})

	return c.service
}
