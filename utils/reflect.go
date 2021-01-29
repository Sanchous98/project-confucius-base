package utils

import (
	"fmt"
	"log"
	"reflect"
)

func CallFunction(i interface{}, name string, args map[string]interface{}) []reflect.Value {
	if !HasFunction(i, name) {
		panic(fmt.Sprintf("Function %s doesn't exist", reflect.TypeOf(i).Name()+name))
	}

	method := reflect.ValueOf(i).MethodByName(name)
	in := make([]reflect.Value, method.Type().NumIn())

	if len(args) != len(in) {
		panic("Invalid number of arguments")
	}

	for k := 0; k < len(in); k++ {
		tt := method.Type().In(k)
		in[k] = reflect.ValueOf(args[tt.String()])
	}

	return method.Call(in)
}

func HasFunction(i interface{}, name string) bool {
	_, ok := reflect.TypeOf(i).MethodByName(name)
	log.Print(reflect.TypeOf(i).MethodByName(name))
	return ok
}

func GetFunctionParamsTypes(i interface{}, name string) []reflect.Type {
	method := reflect.ValueOf(i).MethodByName(name)
	in := make([]reflect.Type, method.Type().NumIn())

	for k := 0; k < len(in); k++ {
		in[k] = method.Type().In(k)
	}

	return in
}
