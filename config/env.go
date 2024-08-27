package config

import (
	"os"
	"reflect"
)

func Scan(env interface{}) map[string]interface{} {
	var mapper = make(map[string]interface{})

	ct := reflect.ValueOf(env).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		val := field.Tag.Get("mapstructure")
		if val != "" {
			mapper[val] = os.Getenv(val)
		}
	}

	return mapper
}
