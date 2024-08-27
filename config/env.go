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

func StringInSlice(value string, slice []string) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}

func MergeStructs(structs ...interface{}) reflect.Type {
	var structFields []reflect.StructField
	var structFieldNames []string

	for _, item := range structs {
		rt := reflect.TypeOf(item)
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			if !StringInSlice(field.Name, structFieldNames) {
				structFields = append(structFields, field)
				structFieldNames = append(structFieldNames, field.Name)
			}
		}
	}

	return reflect.StructOf(structFields)
}
