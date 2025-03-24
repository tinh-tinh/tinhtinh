package transform

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

func ToBool(str interface{}) interface{} {
	typeBool := reflect.TypeOf(str)
	switch typeBool.Kind() {
	case reflect.Bool:
		return str
	case reflect.String:
		val, _ := strconv.ParseBool(str.(string))
		return val
	default:
		panic(fmt.Sprintf("cannot transform bool with type %v, currently only support bool, string", typeBool.Kind()))
	}
}

func ToInt(str interface{}) interface{} {
	typeInt := reflect.TypeOf(str)
	switch typeInt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return str
	case reflect.String:
		val, _ := strconv.Atoi(str.(string))
		return val
	default:
		panic(fmt.Sprintf("cannot transform int with type %v, currently only support int, string", typeInt.Kind()))
	}
}

func ToFloat(str interface{}) interface{} {
	typeFloat := reflect.TypeOf(str)
	switch typeFloat.Kind() {
	case reflect.Float32, reflect.Float64:
		return str
	case reflect.String:
		val, _ := strconv.ParseFloat(str.(string), 64)
		return val
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(str.(int))
	default:
		panic(fmt.Sprintf("cannot transform with type %v, currently only support float, int, string", typeFloat.Kind()))
	}
}

func ToDate(str interface{}) time.Time {
	switch v := str.(type) {
	case time.Time:
		return str.(time.Time)
	case string:
		date, _ := time.Parse("2006-01-02", str.(string))
		return date
	default:
		panic(fmt.Sprintf("cannot transform with type %v, currently only support time, string", v))
	}
}

func ToString(str interface{}) interface{} {
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
