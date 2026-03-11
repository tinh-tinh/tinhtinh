package validator

import (
	"reflect"
	"strconv"
)

// Numeric
func IsInt(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return isInt(value)
}

func isInt(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isInt(value.Index(i)) {
				return false
			}
		}
		return true
	}

	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.String:
		_, err := strconv.Atoi(value.String())
		return err == nil
	default:
		return false
	}
}

func IsFloat(input any) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return isFloat(value)
}

func isFloat(value reflect.Value) bool {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		for i := 0; i < value.Len(); i++ {
			if !isFloat(value.Index(i)) {
				return false
			}
		}
		return true
	}

	switch value.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		_, err := strconv.ParseFloat(value.String(), 64)
		return err == nil
	default:
		return false
	}
}

func IsNumber(str any) bool {
	return IsInt(str) || IsFloat(str)
}
