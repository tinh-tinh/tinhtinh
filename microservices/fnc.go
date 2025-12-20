package microservices

type RpcHandler struct {
	Name    string
	Value   any
	Factory func(c Service, msg Message, reply any) error
}
