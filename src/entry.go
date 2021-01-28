package src

import (
	"github.com/Sanchous98/project-confucius-base/utils"
	"reflect"
	"sync"
)

const constructorName = "Construct"

// containerEntry is a wrapper for services to use binding by abstraction and singletons
type containerEntry struct {
	sync.RWMutex
	service     interface{}
	Abstraction reflect.Type
}

func NewEntry(abstraction reflect.Type, service interface{}) *containerEntry {
	return &containerEntry{
		Abstraction: abstraction,
		service:     service,
	}
}

func (c *containerEntry) Make(container Container) interface{} {
	c.Lock()
	defer c.Unlock()

	if !utils.HasFunction(c.service, constructorName) {
		return c.service
	}

	in := utils.GetFunctionParamsTypes(c.service, constructorName)
	params := make(map[string]interface{}, len(in))

	for _, paramType := range in {
		params[paramType.String()] = container.Get(paramType)
	}

	utils.CallFunction(c.service, constructorName, params)

	return c.service
}
