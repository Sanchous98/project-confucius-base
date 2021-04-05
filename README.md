# Project Confucius framework

Build easily microservice web applications based on [FastHTTP](https://github.com/valyala/fasthttp).

To start you have only to create a "main" function and launch application

```go
package main

import confucius "github.com/Sanchous98/project-confucius-base"

func main() {
	confucius.App().Launch()
}
```

To extend the application you should bind a service to application. Use "inject" tag to bind existent service to your one. 

**Note:** Service must not be bound to application to use, but it every injection will cause instantiating of a service. Otherwise, bound services are singletons
```go
// service.go

package service

type Service struct {
  Log *Log `inject:""`
}

func (s *Service) Constructor() {}

func (s *Service) Destructor() {}
```
```go
// main.go

package main

import confucius "github.com/Sanchous98/project-confucius-base"

func main() {
	confucius.App().Bind(&Service{}).Launch()
}
```

If you want to make a long-living service, you should implement ```Launchable``` interface
```go
type Service struct {
  Log *Log `inject:""`
}

func (s *Service) Launch(chan<- error) {}

func (s *Service) Shutdown(chan<- error) {}
```
**Note:** Launch method is running as a goroutine
"Shutdown" method is called on the app shutdown. Use it to finish service tasks gracefully.
