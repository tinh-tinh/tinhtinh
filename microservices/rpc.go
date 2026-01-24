package microservices

type RpcFactoryFnc func(ctx Ctx) (reply []byte, err error)

type RpcHandler struct {
	Name        string
	Factory     RpcFactoryFnc
	Middlewares []Middleware
}

type RpcHandlers []*RpcHandler
