package utils

import (
	"fmt"
	"reflect"
)

func CallFunction(i interface{}, name string, args []interface{}) []reflect.Value {
	if !HasFunction(i, name) {
		panic(fmt.Sprintf("Function %s doesn't exist", reflect.TypeOf(i).Name()+name))
	}

	method := reflect.ValueOf(i).MethodByName(name)
	in := make([]reflect.Value, method.Type().NumIn())

	if len(args) != len(in) {
		panic("Invalid number of arguments")
	}

	for k := 0; k < len(in); k++ {
		in[k] = reflect.ValueOf(args[k])
	}

	return method.Call(in)
}

func HasFunction(i interface{}, name string) bool {
	_, ok := reflect.TypeOf(i).MethodByName(name)

	if !ok && reflect.Ptr == reflect.TypeOf(i).Kind() {
		_, ok = reflect.ValueOf(i).Type().MethodByName(name)
	}

	return ok
}

func GetMethodParamsTypes(i interface{}, name string) []reflect.Type {
	method := reflect.ValueOf(i).MethodByName(name)
	in := make([]reflect.Type, method.Type().NumIn())

	for k := 0; k < len(in); k++ {
		in[k] = method.Type().In(k)
	}

	return in
}
