package confucius

import (
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/joho/godotenv"
)

type Application struct {
	environment string
	variables   map[string]string
	Container   src.Container
}

func NewApplication(environment string) *Application {
	app := &Application{}
	app.SetEnvironment(environment)
	app.Container = src.NewContainer()

	return app
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
