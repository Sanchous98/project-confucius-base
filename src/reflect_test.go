package src

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type mockServ interface {
	Test()
	Service
}

type extraServ struct{}
type bindServ struct {
	Service *extraServ `inject:""`
}
type fillServ struct {
	Service *bindServ `inject:""`
}
type interfaceServ struct {
	Service *mockServ `inject:""`
}

func (m *fillServ) Test()             {}
func (m *fillServ) Constructor()      {}
func (m *fillServ) Destructor()       {}
func (m *bindServ) Test()             {}
func (m *bindServ) Constructor()      {}
func (m *bindServ) Destructor()       {}
func (m *extraServ) Test()            {}
func (m *extraServ) Constructor()     {}
func (m *extraServ) Destructor()      {}
func (m *interfaceServ) Test()        {}
func (m *interfaceServ) Constructor() {}
func (m *interfaceServ) Destructor()  {}

func TestSingletonBinding(t *testing.T) {
	container := NewContainer()
	serv3 := &extraServ{}
	serv2 := &bindServ{}
	serv := &fillServ{}
	container.Set(serv3)
	container.Set(serv2)

	assert.True(t, container.Has(reflect.TypeOf(serv3)))
	assert.NotPanics(t, func() { fillService(serv, container) })
	assert.NotPanics(t, func() { fillService(serv2, container) })
	assert.NotNil(t, serv.Service)
	assert.NotNil(t, serv2.Service)
}

func TestStructBinding(t *testing.T) {
	container := NewContainer()
	serv2 := &bindServ{}
	serv := &fillServ{}

	assert.NotPanics(t, func() { fillService(serv, container) })
	assert.NotPanics(t, func() { fillService(serv2, container) })
	assert.NotNil(t, serv.Service)
	assert.NotNil(t, serv2.Service)
}

func TestInterfaceBinding(t *testing.T) {
	//var serv mockServ = &bindServ{}
	//container := NewContainer()
	//intServ := &interfaceServ{}
	//container.Set(serv)
	//assert.True(t, container.Has(reflect.TypeOf(serv)))
	//assert.NotPanics(t, func() { fillService(intServ, container) })
	//assert.NotNil(t, intServ.Service)
}
