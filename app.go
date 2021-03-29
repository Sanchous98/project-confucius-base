package confucius

import (
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/joho/godotenv"
	"sync"
)

var (
	application *Application
	once        sync.Once
)

type Application struct {
	environment string
	variables   map[string]string
	container   src.Container
}

func NewApplication(environment string, container src.Container) *Application {
	app := new(Application)
	app.SetEnvironment(environment)
	app.container = container

	return app
}

func App() *Application {
	once.Do(func() {
		application = NewApplication("", src.NewContainer())
	})

	return application
}

func (a *Application) Bind(services ...src.Service) *Application {
	for _, service := range services {
		a.container.Set(service)
	}

	return a
}

func (a Application) Launch() {
	a.container.Launch()
}

func (a *Application) SetEnvironment(name string) {
	var err error
	envFileName := ".env"

	if len(name) > 0 {
		envFileName += "." + name
	}

	a.variables, err = godotenv.Read(envFileName)

	if err != nil {
		panic(err)
	}
}
