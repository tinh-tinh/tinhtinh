package microservices

import "time"

type Service interface {
	Listen()
	Config() Config
}

type ClientProxy interface {
	Config() Config
	Timeout(duration time.Duration) ClientProxy
	Publish(event string, data any, headers ...Header) error
	Send(path string, request any, response any, headers ...Header) error
}
