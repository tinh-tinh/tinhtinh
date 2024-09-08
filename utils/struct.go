package utils

import "reflect"

func GetNameStruct(str interface{}) string {
	name := ""
	if t := reflect.TypeOf(str); t.Kind() == reflect.Ptr {
		name = t.Elem().Name()
	} else {
		name = t.Name()
	}

	return name
}
