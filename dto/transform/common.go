package transform

import (
	"fmt"
	"strconv"
	"time"
)

func StringToBool(str string) bool {
	val, _ := strconv.ParseBool(str)
	return val
}

func ToBool(str interface{}) bool {
	switch v := str.(type) {
	case bool:
		return str.(bool)
	case string:
		val, _ := strconv.ParseBool(str.(string))
		return val
	default:
		panic(fmt.Sprintf("cannot transform bool with type %v", v))
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
		panic(fmt.Sprintf("cannot transform int with type %v", v))
	}
}

func StringToInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func StringToInt(str string) int {
	val, _ := strconv.Atoi(str)
	return val
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
		panic(fmt.Sprintf("cannot transform with type %v", v))
	}
}

func StringToDate(str string) time.Time {
	date, _ := time.Parse("2006-01-02", str)
	return date
}

func StringToTimeDuration(str string) time.Duration {
	val, _ := time.ParseDuration(str)
	return val
}
