package core

type Handler func(ctx Ctx)

type Controller interface {
	Guard(...Guard) *Controller
	Pipe(...Middleware) *Controller
	Get(string, Handler)
	Post(string, Handler)
	Patch(string, Handler)
	Put(string, Handler)
	Delete(string, Handler)
}
