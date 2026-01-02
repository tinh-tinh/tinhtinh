package validator

import (
	"reflect"
	"time"
)

// Date time
func IsDate(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return isDate(value)
}

func isDate(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isDate(value.Index(i)) {
				return false
			}
		}
		return true
	}

	switch value.Kind() {
	case reflect.String:
		_, err := time.Parse("2006-01-02", value.String())
		return err == nil
	case reflect.Struct:
		if value.Type() == reflect.TypeOf(time.Time{}) {
			return true
		}
		return false
	default:
		return false
	}
}
