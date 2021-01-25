package main

import (
	confucius "github.com/Sanchous98/project-confucius-base"
	"github.com/Sanchous98/project-confucius-base/service/graphql"
	"github.com/Sanchous98/project-confucius-base/service/static"
	"github.com/Sanchous98/project-confucius-base/service/web"
	"github.com/Sanchous98/project-confucius-base/utils"
	"reflect"
)

var application = confucius.NewApplication("example")

func main() {
	webService := web.Web{}
	staticService := static.Static{}
	graphqlService := graphql.GraphQL{}
	application.SetErrorHandler(utils.AssertError())
	application.Container.Set(reflect.TypeOf(webService), &webService, true)
	application.Container.Set(reflect.TypeOf(staticService), &staticService, true)
	application.Container.Set(reflect.TypeOf(graphqlService), &graphqlService, true)
	application.Container.Launch()
}
