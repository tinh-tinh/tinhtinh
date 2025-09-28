package microservices

import "time"

type Service interface {
	Listen()
	Config() Config
}

type ClientProxy interface {
	Config() Config
	Timeout(duration time.Duration) ClientProxy
	Publish(event string, data interface{}, headers ...Header) error
}
