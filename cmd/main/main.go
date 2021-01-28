package main

import (
	confucius "github.com/Sanchous98/project-confucius-base"
	"github.com/Sanchous98/project-confucius-base/service/graphql"
	"github.com/Sanchous98/project-confucius-base/service/static"
	"github.com/Sanchous98/project-confucius-base/service/web"
)

func main() {
	webService := web.Web{}
	staticService := static.Static{}
	graphqlService := graphql.GraphQL{}
	confucius.App().Container.Set(&webService)
	confucius.App().Container.Set(&staticService)
	confucius.App().Container.Set(&graphqlService)
	confucius.App().Container.Launch()
}
