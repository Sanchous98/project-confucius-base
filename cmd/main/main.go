package main

import (
	confucius "github.com/Sanchous98/project-confucius-base"
	"github.com/Sanchous98/project-confucius-base/service/graphql"
	"github.com/Sanchous98/project-confucius-base/service/static"
	"github.com/Sanchous98/project-confucius-base/service/web"
)

func main() {
	confucius.App().Container.Set(&web.Web{})
	confucius.App().Container.Set(&static.Static{})
	confucius.App().Container.Set(&graphql.GraphQL{})
	confucius.App().Container.Launch()
}
