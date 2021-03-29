package main

import (
	confucius "github.com/Sanchous98/project-confucius-base"
	"github.com/Sanchous98/project-confucius-base/stdlib"
)

func main() {
	confucius.
		App().
		Bind(&stdlib.Web{}, &stdlib.Static{}, &stdlib.GraphQL{}).
		Launch()
}
