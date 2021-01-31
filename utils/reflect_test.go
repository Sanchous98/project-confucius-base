package utils

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"unsafe"
)

type Mock struct{}

func (m *Mock) ByPointer(integer int, char string) unsafe.Pointer {
	return unsafe.Pointer(m)
}

func (m *Mock) Returning(integer int, integer2 *int) (int, *int) {
	return integer, integer2
}

func (m Mock) ByValue() unsafe.Pointer {
	return unsafe.Pointer(&m)
}

func TestHasFunction(t *testing.T) {
	assert.True(t, HasFunction(&Mock{}, "ByPointer"))
	assert.True(t, HasFunction(Mock{}, "ByValue"))
}

func TestGetMethodParamsTypes(t *testing.T) {
	params := GetMethodParamsTypes(&Mock{}, "ByPointer")
	assert.Equal(t, reflect.Int, params[0].Kind())
	assert.Equal(t, reflect.String, params[1].Kind())
}

func TestCallFunction(t *testing.T) {
	i := 7
	returns := CallFunction(&Mock{}, "Returning", []interface{}{6, &i})
	assert.Equal(t, 6, returns[0].Interface())
	assert.Equal(t, 7, returns[1].Elem().Interface())
}
