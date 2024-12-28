package microservices

import (
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/logger"
)

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

type ConnectOptions struct {
	Addr         string
	Serializer   core.Encode
	Deserializer core.Decode
	Timeout      time.Duration
	ErrorHandler ErrorHandler
	Logger       *logger.Logger
}

type ClientProxy interface {
	SetHeaders(key string, value string) ClientProxy
	GetHeaders(key string) string
	Send(event string, data interface{}) error
	Publish(event string, data interface{}) error
}
