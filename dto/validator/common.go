package validator

import (
	"fmt"
	"reflect"
)

const MiB = 1 << 20 // 1 MiB

func typeof(v any) string {
	return fmt.Sprintf("%T", v)
}

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	return isEmpty(val)
}

func isEmpty(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Pointer:
		return v.IsNil()
	case reflect.Array, reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	default:
		return v.IsZero()
	}
}

func MinLength(input any, min int) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return minLength(value, min)
}

func minLength(value reflect.Value, min int) bool {
	if !value.IsValid() {
		return false
	}

	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		return value.Len() >= min
	}

	if value.Kind() == reflect.String {
		return len(value.String()) >= min
	}

	return false
}

func MaxLength(input any, max int) bool {
	if input == nil {
		return false
	}

	value := reflect.ValueOf(input)
	return maxLength(value, max)
}

func maxLength(value reflect.Value, max int) bool {
	if !value.IsValid() {
		return false
	}

	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		return value.Len() <= max
	}

	if value.Kind() == reflect.String {
		return len(value.String()) <= max
	}

	return false
}
