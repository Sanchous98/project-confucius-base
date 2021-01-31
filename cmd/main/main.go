package main

import (
	confucius "github.com/Sanchous98/project-confucius-base"
	"github.com/Sanchous98/project-confucius-base/lib"
)

func main() {
	confucius.App().Container.Set(&lib.Web{})
	confucius.App().Container.Set(&lib.Static{})
	confucius.App().Container.Set(&lib.GraphQL{})
	confucius.App().Container.Launch()
}
