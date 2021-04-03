package confucius

import (
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/Sanchous98/project-confucius-base/stdlib"
	"github.com/joho/godotenv"
	"os"
	"sync"
)

var (
	app  *application
	once sync.Once
)

type application struct {
	environment string
	variables   map[string]string
	container   src.Container
}

func NewApplication(environment string, container src.Container) *application {
	app := new(application)
	app.SetEnvironment(environment)
	app.container = container

	return app
}

func App() *application {
	once.Do(func() {
		app = NewApplication("", src.NewContainer())
		app.bootstrap()
	})

	return app
}

func (a *application) bootstrap() {
	a.Bind(&stdlib.Web{}, &stdlib.Static{}).Launch()
}

func (a *application) Bind(services ...src.Service) *application {
	for _, service := range services {
		a.container.Set(service)
	}

	return a
}

func (a application) Launch() {
	a.container.Launch()
}

func (a *application) SetEnvironment(name string) {
	var err error
	envFileName := ".env"

	if len(name) > 0 {
		envFileName += "." + name
	}

	a.variables, err = godotenv.Read(envFileName)

	if err != nil {
		_, err = os.Create(envFileName)

		if err != nil {
			panic(err)
		}
	}
}

func Getenv(name string) string {
	return App().variables[name]
}

func Setenv(name, value string) {
	App().variables[name] = value
}
