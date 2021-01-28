package src

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
	"unsafe"
)

var container Container = NewContainer()

type MockService struct {
	sync.Mutex
	mockData int
}

func (m *MockService) Make(Container) *MockService {
	return m
}

// TestSingletonBinding checks whether bound services are not re-instantiated, when getting them from container.
// Pointer values must not change after second call of getter.
func TestSingletonBinding(t *testing.T) {
	mock := new(MockService)
	container.Set(reflect.TypeOf(mock), mock)

	assert.True(t, container.Has(reflect.TypeOf(mock)))

	serviceFirst := container.Get(reflect.TypeOf(mock)).(*MockService)
	serviceSecond := container.Get(reflect.TypeOf(mock)).(*MockService)

	assert.Equal(t, uintptr(unsafe.Pointer(serviceFirst)), uintptr(unsafe.Pointer(serviceSecond)))
}

// TestMultitonBinding checks whether bound services are re-instantiated, when getting them from container
// Pointer values must not change after second call of getter
func TestMultitonBinding(t *testing.T) {
	mock := new(MockService)
	container.Set(reflect.TypeOf(mock), mock)

	assert.True(t, container.Has(reflect.TypeOf(mock)))

	serviceFirst := container.Get(reflect.TypeOf(mock)).(*MockService)
	serviceSecond := container.Get(reflect.TypeOf(mock)).(*MockService)

	assert.NotEqual(t, uintptr(unsafe.Pointer(serviceFirst)), uintptr(unsafe.Pointer(serviceSecond)))
}

func TestMute(t *testing.T) {
	mock := new(MockService)
	container.Set(reflect.TypeOf(mock), mock)
	mock2 := container.Get(reflect.TypeOf(mock)).(*MockService)

	assert.True(t, container.Has(reflect.TypeOf(mock)))
	assert.Equal(t, uintptr(unsafe.Pointer(mock)), uintptr(unsafe.Pointer(mock2)))

	go modifyMock(mock, 1)
	assert.Equal(t, mock.mockData, 1)

	assert.Equal(t, mock.mockData, 2)
}

func modifyMock(mock *MockService, value int) {
	mock.Lock()
	mock.mockData = value
	mock.Unlock()
}
