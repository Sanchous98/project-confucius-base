package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

type Mock struct{}

func (m *Mock) ByPointer() unsafe.Pointer {
	return unsafe.Pointer(m)
}

func (m Mock) ByValue() unsafe.Pointer {
	return unsafe.Pointer(&m)
}

func TestHasFunction(t *testing.T) {
	assert.True(t, HasFunction(Mock{}, "ByPointer"))
}

func TestGetFunctionParamsTypes(t *testing.T) {

}

func TestCallFunction(t *testing.T) {

}
