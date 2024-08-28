package utils

import (
	"strconv"
	"time"
)

func StringToInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return val
}

func StringToBool(str string) bool {
	val, err := strconv.ParseBool(str)
	if err != nil {
		panic(err)
	}
	return val
}

func StringToTimeDuration(str string) time.Duration {
	val, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}
	return val
}
