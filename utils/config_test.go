package utils

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

type mockConfig struct {
	TestConfig string        `yaml:"test_config"`
	ArrayMock  []interface{} `yaml:"array_mock"`
}

func (m *mockConfig) Unmarshall() error {
	return nil
}

func TestHydrateConfig(t *testing.T) {
	config, err := HydrateConfig(&mockConfig{}, "/../test/config/mock.yml", yaml.Unmarshal)

	assert.NoError(t, err)

	mock := config.(*mockConfig)

	assert.Equal(t, mock.TestConfig, "mock", "mock.TestConfig expected \"mock\", \"%s\" given", mock.TestConfig)
	assert.Equal(t, mock.ArrayMock[0].(string), "test1", "mock.ArrayMock has got invalid elements. Expected \"test1\"")
	assert.Equal(t, mock.ArrayMock[1].(string), "test2", "mock.ArrayMock has got invalid elements. Expected \"test2\"")
}
