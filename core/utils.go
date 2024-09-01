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
