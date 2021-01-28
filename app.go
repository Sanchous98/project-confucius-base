package confucius

import (
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/joho/godotenv"
	"os"
	"sync"
)

var (
	application *Application
	once        sync.Once
)

type Application struct {
	environment string
	variables   map[string]string
	Container   src.Container
}

func NewApplication(environment string, container src.Container) *Application {
	app := new(Application)
	app.SetEnvironment(environment)
	app.Container = container

	return app
}

func App() *Application {
	once.Do(func() {
		application = NewApplication(os.Getenv("APP_ENV"), src.NewContainer())
	})

	return application
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
