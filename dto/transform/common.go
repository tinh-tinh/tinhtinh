package transform

import (
	"fmt"
	"strconv"
	"time"
)

func ToBool(str interface{}) bool {
	switch v := str.(type) {
	case bool:
		return str.(bool)
	case string:
		val, _ := strconv.ParseBool(str.(string))
		return val
	default:
		panic(fmt.Sprintf("cannot transform bool with type %v, currently only support bool, string", v))
	}
}

func ToInt(str interface{}) interface{} {
	switch v := str.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return str
	case string:
		val, _ := strconv.Atoi(str.(string))
		return val
	default:
		panic(fmt.Sprintf("cannot transform int with type %v, currently only support int, string", v))
	}
}

func ToFloat(str interface{}) interface{} {
	switch v := str.(type) {
	case float32:
		return str.(float32)
	case float64:
		return str.(float64)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return float64(str.(int))
	case string:
		val, _ := strconv.ParseFloat(str.(string), 64)
		return val
	default:
		panic(fmt.Sprintf("cannot transform with type %v, currently only support float, int, string", v))
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

func ToString(str interface{}) string {
	switch v := str.(type) {
	case string:
		return str.(string)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return strconv.Itoa(int(str.(int)))
	case float32:
		return strconv.FormatFloat(float64(str.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(str.(float64), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(str.(bool))
	case time.Time:
		return str.(time.Time).String()
	default:
		panic(fmt.Sprintf("cannot transform with type %v, currently only support string, int, float, bool, time", v))
	}
}
