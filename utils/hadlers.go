package utils

import (
	"fmt"
	"log"
)

type ErrorHandler interface {
	Catch(error)
	OnError(func(error))
}

// Global error handler
type assert struct {
	err      chan error
	callback func(error)
}

func AssertError() *assert {
	a := &assert{make(chan error), defaultHandler}

	go func() {
		for {
			select {
			case errors := <-a.err:
				a.callback(errors)
			}
		}
	}()

	return a
}

func (a *assert) OnError(callback func(error)) {
	a.callback = callback
}

func (a *assert) Catch(err error) {
	a.err <- err
}

func defaultHandler(err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%e", err))
	}
}
