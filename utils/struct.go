package utils

import (
	"reflect"
	"runtime"
)

func GetNameStruct(str interface{}) string {
	name := ""
	if t := reflect.TypeOf(str); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	return name
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
