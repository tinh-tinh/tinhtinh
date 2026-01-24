package microservices

type RPCGateway struct {
	Handler func(ctx Ctx) ([]byte, error)
}

func (r *RPCGateway) Call(ctx Ctx) ([]byte, error) {
	return r.Handler(ctx)
}
