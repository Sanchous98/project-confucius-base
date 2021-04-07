package project_confucius_base

import (
	"fmt"
	"github.com/Sanchous98/project-confucius-base/src"
	"github.com/Sanchous98/project-confucius-base/stdlib"
	"github.com/Sanchous98/project-confucius-base/utils"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const appConfigPath = "config/app.yaml"

var (
	app  *application
	once sync.Once
)

type catcher func(chan error, chan os.Signal, *application)

type application struct {
	environment  string
	variables    map[string]string
	container    src.Container
	errorCatcher catcher
	config       *appConfig
	logger       stdlib.Logger
}

type appConfig struct {
	AlwaysRestart bool `yaml:"always_restart"`
}

func (a *appConfig) Unmarshall() error {
	return utils.Unmarshall(a, appConfigPath, yaml.Unmarshal)
}

func NewApplication(environment string, container src.Container) *application {
	app := new(application)
	app.SetEnvironment(environment)
	app.container = container

	return app
}

func (a *application) SetErrorCatcher(catcher catcher) {
	a.errorCatcher = catcher
}

func App() *application {
	once.Do(func() {
		app = NewApplication("", src.NewContainer())
		app.bootstrap()
	})

	return app
}

func (a *application) bootstrap() {
	a.errorCatcher = defaultCatcher
	app.logger = &stdlib.Log{}
	app.config = new(appConfig)
	err := app.config.Unmarshall()

	if err != nil {
		app.logger.Emergency(err)
	}

	a.Bind(&stdlib.Web{}, &stdlib.Static{}, app.logger, &stdlib.Metrics{})
}

func (a *application) Bind(services ...src.Service) *application {
	for _, service := range services {
		a.container.Set(service)
	}

	return a
}

func (a *application) Launch(alwaysRestart bool) {
	err := make(chan error)
	osSignals := make(chan os.Signal)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	for _, service := range a.container.GetServices() {
		switch launchable := service.(type) {
		case src.Launcher:
			go func() {
				defer defaultRecover(a, alwaysRestart)
				launchable.Launch()
			}()
		}
	}

	a.errorCatcher(err, osSignals, a)
}

func (a *application) Shutdown() {
	for _, service := range a.container.GetServices() {
		switch launchable := service.(type) {
		case src.Launcher:
			go launchable.Shutdown()
		}
	}
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

func defaultRecover(a *application, restart bool) {
	if r := recover(); r != nil {
		log.Print(r)
		a.Shutdown()

		if restart {
			fmt.Println("Restarting app")
			time.Sleep(time.Second)
			a.Launch(restart)
		} else {
			os.Exit(1)
		}
	}
}

func defaultCatcher(errors chan error, signals chan os.Signal, a *application) {
	select {
	case <-signals:
		log.Print("\nShutdown signal received.")
		a.Shutdown()
		log.Print("\nServer gracefully stopped.\n")
		os.Exit(0)
	case err := <-errors:
		if err != nil {
			a.Shutdown()
			log.Fatal(err)
		}
	}
}
