package core

type CallHandler func(data Map) Map

type Interceptor func(ctx Ctx) CallHandler

func (c *DynamicController) Interceptor(interceptor Interceptor) Controller {
	c.interceptor = interceptor
	return c
}
