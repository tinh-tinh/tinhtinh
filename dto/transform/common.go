package transform

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ToBool(val any) any {
	if val == nil {
		panic(fmt.Errorf("cannot convert nil to bool"))
	}

	v := reflect.ValueOf(val)

	// Handle array or slice
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		length := v.Len()
		if length == 0 {
			return val
		}

		// Find common type in the slice
		var commonType reflect.Type
		results := make([]any, length)

		for i := range length {
			converted := ToBool(v.Index(i).Interface())
			results[i] = converted

			if commonType == nil {
				commonType = reflect.TypeOf(converted)
			} else if commonType != reflect.TypeOf(converted) {
				panic(fmt.Errorf("inconsistent types in array: found %v and %v", commonType, reflect.TypeOf(converted)))
			}
		}

		// Convert to slice
		outputSlice := reflect.MakeSlice(reflect.SliceOf(commonType), length, length)
		for i := range length {
			outputSlice.Index(i).Set(reflect.ValueOf(results[i]))
		}
		return outputSlice.Interface()
	}

	// Handle single value
	typeBool := reflect.TypeOf(val)
	switch typeBool.Kind() {
	case reflect.Bool:
		return val
	case reflect.String:
		val, err := strconv.ParseBool(val.(string))
		if err != nil {
			panic(fmt.Sprintf("cannot transform bool with type %v, currently only support bool, string", typeBool.Kind()))
		}
		return val
	default:
		panic(fmt.Sprintf("cannot transform bool with type %v, currently only support bool, string", typeBool.Kind()))
	}
}

func ToInt(val any) any {
	if val == nil {
		panic(fmt.Errorf("cannot convert nil to int"))
	}

	v := reflect.ValueOf(val)

	// Handle array or slice
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		length := v.Len()
		if length == 0 {
			return val
		}

		// Find common type in the slice
		var commonType reflect.Type
		results := make([]any, length)

		for i := range length {
			converted := ToInt(v.Index(i).Interface())
			results[i] = converted

			if commonType == nil {
				commonType = reflect.TypeOf(converted)
			} else if commonType != reflect.TypeOf(converted) {
				panic(fmt.Errorf("inconsistent types in array: found %v and %v", commonType, reflect.TypeOf(converted)))
			}
		}

		// Convert to slice
		outputSlice := reflect.MakeSlice(reflect.SliceOf(commonType), length, length)
		for i := range length {
			outputSlice.Index(i).Set(reflect.ValueOf(results[i]))
		}
		return outputSlice.Interface()
	}

	// Handle one element
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val

	case reflect.String:
		num, err := strconv.Atoi(val.(string))
		if err != nil {
			panic(fmt.Errorf("cannot convert string '%v' to int: %v", val, err))
		}
		return num

	default:
		panic(fmt.Errorf("cannot convert type %v to int, only support int, uint, string, and array/slice of them", v.Kind()))
	}
}

func ToFloat(str any) any {
	if str == nil {
		panic(fmt.Errorf("cannot convert nil to float"))
	}

	v := reflect.ValueOf(str)

	// Handle array or slice
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		length := v.Len()
		if length == 0 {
			return str
		}

		// Find common type in the slice
		var commonType reflect.Type
		results := make([]any, length)

		for i := range length {
			converted := ToFloat(v.Index(i).Interface())
			results[i] = converted

			if commonType == nil {
				commonType = reflect.TypeOf(converted)
			} else if commonType != reflect.TypeOf(converted) {
				panic(fmt.Errorf("inconsistent types in array: found %v and %v", commonType, reflect.TypeOf(converted)))
			}
		}

		// Convert to slice
		outputSlice := reflect.MakeSlice(reflect.SliceOf(commonType), length, length)
		for i := range length {
			outputSlice.Index(i).Set(reflect.ValueOf(results[i]))
		}
		return outputSlice.Interface()
	}

	// Handle one element
	typeFloat := reflect.TypeOf(str)
	switch typeFloat.Kind() {
	case reflect.Float32, reflect.Float64:
		return str
	case reflect.String:
		val, err := strconv.ParseFloat(str.(string), 64)
		if err != nil {
			panic(fmt.Sprintf("cannot transform with type %v, currently only support float, int, string", typeFloat.Kind()))
		}
		return val
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(str.(int))
	default:
		panic(fmt.Sprintf("cannot transform with type %v, currently only support float, int, string", typeFloat.Kind()))
	}
}

func ToDate(str any) any {
	if str == nil {
		panic(fmt.Errorf("cannot convert nil to date"))
	}

	v := reflect.ValueOf(str)

	// Handle array or slice
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		length := v.Len()
		if length == 0 {
			return str
		}

		// Find common type in the slice
		results := make([]time.Time, length)

		for i := range length {
			converted := ToDate(v.Index(i).Interface())
			results[i] = converted.(time.Time)
		}

		return results
	}
	switch v := str.(type) {
	case time.Time:
		return str.(time.Time)
	case string:
		date, err := time.Parse("2006-01-02", str.(string))
		if err != nil {
			panic(fmt.Sprintf("cannot transform with type %v, currently only support time, string", v))
		}
		return date
	default:
		panic(fmt.Sprintf("cannot transform with type %v, currently only support time, string", v))
	}
}

func ToString(str any) any {
	if str == nil {
		return ""
	}

	v := reflect.ValueOf(str)

	// Handle array or slice
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		length := v.Len()
		if length == 0 {
			return str
		}

		// Find common type in the slice
		var commonType reflect.Type
		results := make([]any, length)

		for i := range length {
			converted := ToString(v.Index(i).Interface())
			results[i] = converted

			if commonType == nil {
				commonType = reflect.TypeOf(converted)
			} else if commonType != reflect.TypeOf(converted) {
				panic(fmt.Errorf("inconsistent types in array: found %v and %v", commonType, reflect.TypeOf(converted)))
			}
		}

		// Convert to slice
		outputSlice := reflect.MakeSlice(reflect.SliceOf(commonType), length, length)
		for i := range length {
			outputSlice.Index(i).Set(reflect.ValueOf(results[i]))
		}

		return outputSlice.Interface()
	}

	typeStr := reflect.TypeOf(str)
	switch typeStr.Kind() {
	case reflect.String:
		return str
	case reflect.Int:
		return strconv.Itoa(str.(int))
	case reflect.Int8:
		return strconv.Itoa(int(str.(int8)))
	case reflect.Int16:
		return strconv.Itoa(int(str.(int16)))
	case reflect.Int32:
		return strconv.Itoa(int(str.(int32)))
	case reflect.Int64:
		return strconv.Itoa(int(str.(int64)))
	case reflect.Uint:
		return strconv.Itoa(int(str.(uint)))
	case reflect.Uint8:
		return strconv.Itoa(int(str.(uint8)))
	case reflect.Uint16:
		return strconv.Itoa(int(str.(uint16)))
	case reflect.Uint32:
		return strconv.Itoa(int(str.(uint32)))
	case reflect.Uint64:
		return strconv.Itoa(int(str.(uint64)))
	case reflect.Float32:
		return strconv.FormatFloat(float64(str.(float32)), 'f', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(str.(float64), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(str.(bool))
	case reflect.Struct:
		if typeStr == reflect.TypeOf(time.Time{}) {
			return str.(time.Time).String()
		}
		panic(fmt.Sprintf("cannot transform with type %v, currently only support string, int, float, bool, time", typeStr.Kind()))
	default:
		panic(fmt.Sprintf("cannot transform with type %v, currently only support string, int, float, bool, time", typeStr.Kind()))
	}
}
