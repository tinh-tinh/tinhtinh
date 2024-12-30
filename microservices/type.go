package microservices

type EventType string

const (
	RPC    EventType = "rpc"
	PubSub EventType = "pubsub"
)

type Message struct {
	Type    EventType         `json:"type"`
	Event   string            `json:"event"`
	Headers map[string]string `json:"headers"`
	Data    interface{}       `json:"data"`
}

type Service interface {
	Listen()
	Serializer(v interface{}) ([]byte, error)
	Deserializer(data []byte, v interface{}) error
	ErrorHandler(err error)
}

type ClientProxy interface {
	Send(event string, data interface{}, headers ...Header) error
	Publish(event string, data interface{}, headers ...Header) error
}

type Options struct {
	Config
	Addr string
}
