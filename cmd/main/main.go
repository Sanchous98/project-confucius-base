package main

import (
	confucius "github.com/Sanchous98/project-confucius-base"
	"github.com/Sanchous98/project-confucius-base/service/graphql"
	"github.com/Sanchous98/project-confucius-base/service/static"
	"github.com/Sanchous98/project-confucius-base/service/web"
	"log"
	"reflect"
	"unsafe"
)

func main() {
	webService := web.Web{}
	staticService := static.Static{}
	graphqlService := graphql.GraphQL{}
	log.Print(unsafe.Pointer(&webService))
	confucius.App().Container.Set(reflect.TypeOf(webService), &webService)
	confucius.App().Container.Set(reflect.TypeOf(staticService), &staticService)
	confucius.App().Container.Set(reflect.TypeOf(graphqlService), &graphqlService)
	confucius.App().Container.Launch()
}
