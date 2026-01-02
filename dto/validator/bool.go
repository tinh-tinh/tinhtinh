package validator

import (
	"reflect"
	"strconv"
)

// Boolean
func IsBool(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return isBool(value)
}

func isBool(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isBool(value.Index(i)) {
				return false
			}
		}
		return true
	}

	switch value.Kind() {
	case reflect.Bool:
		return true
	case reflect.String:
		_, err := strconv.ParseBool(value.String())
		return err == nil
	default:
		return false
	}
}
