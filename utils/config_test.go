package utils

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"testing"
)

const testConfigPath = "testdata/test_config.yaml"

type mockConfig struct {
	TestConfig string        `yaml:"test_config"`
	ArrayMock  []interface{} `yaml:"array_mock"`
}

func (m *mockConfig) Unmarshall() error {
	return nil
}

func TestHydrateConfig(t *testing.T) {
	var mock mockConfig
	err := Unmarshall(&mock, testConfigPath, yaml.Unmarshal)

	assert.NoError(t, err)
	assert.Equal(t, "mock", mock.TestConfig, "Mock.TestConfig expected \"Mock\", \"%s\" given", mock.TestConfig)
	assert.Equal(t, "test1", mock.ArrayMock[0].(string), "Mock.ArrayMock has got invalid elements. Expected \"test1\"")
	assert.Equal(t, "test2", mock.ArrayMock[1].(string), "Mock.ArrayMock has got invalid elements. Expected \"test2\"")
}
