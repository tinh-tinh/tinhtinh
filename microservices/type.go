package microservices

import "time"

type Service interface {
	Listen()
	Config() Config
}

type ClientProxy interface {
	Config() Config
	Timeout(duration time.Duration) ClientProxy
	Emit(event string, message Message) error
	Send(event string, data interface{}, headers ...Header) error
	Publish(event string, data interface{}, headers ...Header) error
}
