package core

import (
	"strconv"
	"strings"
)

type Route struct {
	Method string
	Path   string
}

func ParseRoute(url string) Route {
	route := strings.Split(url, " ")

	var path string
	for i := 1; i < len(route); i++ {
		path += IfSlashPrefixString(route[i])
	}

	return Route{
		Method: route[0],
		Path:   path,
	}
}

func (r *Route) SetPrefix(prefix string) {
	r.Path = IfSlashPrefixString(prefix) + r.Path
}

func (r *Route) GetPath() string {
	return r.Method + " " + r.Path
}

func IfSlashPrefixString(s string) string {
	if s == "" {
		return s
	}
	s = strings.TrimSuffix(s, "/")
	if strings.HasPrefix(s, "/") {
		return ToFormat(s)
	}
	return "/" + ToFormat(s)
}

func ToFormat(s string) string {
	result := strings.ToLower(s)
	return strings.ReplaceAll(result, " ", "")
}

func IntToString(a int) string {
	s := strconv.Itoa(a)
	return s
}

const (
	Red     = "\x1b[31m"
	Green   = "\x1b[32m"
	Yellow  = "\x1b[33m"
	Blue    = "\x1b[34m"
	Magenta = "\x1b[35m"
	Cyan    = "\x1b[36m"
	Reset   = "\x1b[0m"
	Gray    = "\033[37m"
	White   = "\033[97m"
)
