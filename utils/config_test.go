package utils

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

var mock = []byte("test_config: mock\narray_mock:\n  - test1\n  - test2\n  - sub_array:\n      - test3\n      - test4")

type mockConfig struct {
	TestConfig string        `yaml:"test_config"`
	ArrayMock  []interface{} `yaml:"array_mock"`
}

func (m *mockConfig) Unmarshall() error {
	return nil
}

func TestHydrateConfig(t *testing.T) {
	config, err := HydrateConfig(&mockConfig{}, mock, yaml.Unmarshal)

	assert.NoError(t, err)

	mock := config.(*mockConfig)

	assert.Equal(t, mock.TestConfig, "Mock", "Mock.TestConfig expected \"Mock\", \"%s\" given", mock.TestConfig)
	assert.Equal(t, mock.ArrayMock[0].(string), "test1", "Mock.ArrayMock has got invalid elements. Expected \"test1\"")
	assert.Equal(t, mock.ArrayMock[1].(string), "test2", "Mock.ArrayMock has got invalid elements. Expected \"test2\"")
}
