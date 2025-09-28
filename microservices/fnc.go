package microservices

type ReqFnc func(event string, data interface{}, headers ...Header) error
