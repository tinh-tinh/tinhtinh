package core

import (
	"github.com/tinh-tinh/tinhtinh/v2/common"
)

const MethodAll = "ALL"

type RoutesPath struct {
	Path   string
	Method string
}

func (c *Consumer) Include(includes ...RoutesPath) *Consumer {
	c.includes = append(c.includes, includes...)
	return c
}

func (c *Consumer) Exclude(excludes ...RoutesPath) *Consumer {
	c.excludes = append(c.excludes, excludes...)
	return c
}

type Consumer struct {
	middlewares []Middleware
	includes    []RoutesPath
	excludes    []RoutesPath
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Apply(middlewares ...Middleware) *Consumer {
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

func (m *DynamicModule) Consumer(consumer *Consumer) Module {
	effectRoutes := []*Router{}
	for _, i := range consumer.includes {
		if i.Path == "*" && i.Method == MethodAll {
			effectRoutes = m.Routers
		} else if i.Path == "*" {
			effectRoutes = common.Filter(m.Routers, func(r *Router) bool {
				return r.Method == i.Method
			})
		} else if i.Method == MethodAll {
			effectRoutes = common.Filter(m.Routers, func(r *Router) bool {
				route := ParseRoute(" " + r.Path)
				route.SetPrefix(r.Name)
				return route.Path == i.Path
			})
		} else {
			effectRoutes = common.Filter(m.Routers, func(r *Router) bool {
				route := ParseRoute(" " + r.Path)
				route.SetPrefix(r.Name)
				return r.Method == i.Method && route.Path == i.Path
			})
		}
	}
	if len(consumer.includes) == 0 {
		effectRoutes = m.Routers
	}

	for _, e := range consumer.excludes {
		if e.Path == "*" && e.Method == MethodAll {
			effectRoutes = common.Remove(effectRoutes, func(r *Router) bool {
				return true
			})
		} else if e.Path == "*" {
			effectRoutes = common.Remove(effectRoutes, func(r *Router) bool {
				return r.Method == e.Method
			})
		} else if e.Method == MethodAll {
			effectRoutes = common.Remove(effectRoutes, func(r *Router) bool {
				route := ParseRoute(" " + r.Path)
				route.SetPrefix(r.Name)
				return route.Path == e.Path
			})
		} else {
			effectRoutes = common.Remove(effectRoutes, func(r *Router) bool {
				route := ParseRoute(" " + r.Path)
				route.SetPrefix(r.Name)
				return r.Method == e.Method && route.Path == e.Path
			})
		}
	}

	for _, r := range effectRoutes {
		r.Middlewares = append(consumer.middlewares, r.Middlewares...)
	}

	return m
}
