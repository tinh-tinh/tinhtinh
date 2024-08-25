package api

import "strconv"

func IntToString(a int) string {
	s := strconv.Itoa(a)
	return s
}
