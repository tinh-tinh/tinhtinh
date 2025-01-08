package microservices

type Service interface {
	Listen()
	Serializer(v interface{}) ([]byte, error)
	Deserializer(data []byte, v interface{}) error
	ErrorHandler(err error)
}

type ClientProxy interface {
	Emit(event string, message Message) error
	Headers() Header
	Serializer(v interface{}) ([]byte, error)
	Deserializer(data []byte, v interface{}) error
	Send(event string, data interface{}, headers ...Header) error
	Publish(event string, data interface{}, headers ...Header) error
}
