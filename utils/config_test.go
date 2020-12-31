package utils

import (
	"gopkg.in/yaml.v2"
	"os"
	"testing"
)

type mockConfig struct {
	TestConfig string        `yaml:"test_config"`
	ArrayMock  []interface{} `yaml:"array_mock"`
}

func (m *mockConfig) HydrateConfig() error {
	return nil
}

func TestHydrateConfig(t *testing.T) {
	workDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	config, err := HydrateConfig(&mockConfig{}, workDir+"/test/config/mock.yml", yaml.Unmarshal)

	if err != nil {
		t.Fatal(err)
	}

	mock := config.(*mockConfig)

	if mock.TestConfig != "mock" {
		t.Fatalf("mock.TestConfig expected \"mock\", \"%s\" given", mock.TestConfig)
	}

	if mock.ArrayMock[0].(string) != "test1" || mock.ArrayMock[1].(string) != "test2" {
		t.Fatal("mock.ArrayMock has got invalid elements. Expected \"test1\" and \"test2\"")
	}
}
