package core

type CallHandler func(data Map) Map

type Interceptor func(ctx *Ctx) CallHandler

func (c *DynamicController) Interceptor(interceptor Interceptor) *DynamicController {
	c.interceptor = interceptor
	return c
}
