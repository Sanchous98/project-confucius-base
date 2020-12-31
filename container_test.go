package confucius

import (
	"github.com/gorilla/mux"
	"reflect"
	"testing"
)

var container Container

type MockService struct{}

func (m MockService) Serve(handler *mux.Router) error {
	return nil
}
func (m MockService) Stop() {}
func (m MockService) Init() error {
	return nil
}

func TestRegisterService(t *testing.T) {
	container = NewContainer(&Config{})
	container.Set("mock_service", MockService{})

	if !container.Has("mock_service") {
		t.Fatal("Failed to register a service")
	} else {
		service, _ := container.Get("mock_service")

		if reflect.TypeOf(service) == reflect.TypeOf(MockService{}) {
			t.Fatal("Invalid type of registered service")
		}
	}
}

func TestInitServices(t *testing.T) {
	TestRegisterService(t)
	err := container.Init()

	if err != nil {
		t.Fatal(err)
	}

	_, status := container.Get("mock_service")

	if status != Ok {
		t.Fatal("Mock service not initialized")
	}
}
