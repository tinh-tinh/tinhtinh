package microservices

import "github.com/tinh-tinh/tinhtinh/v2/core"

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type Service interface {
	Listen()
}

type ConnectOptions struct {
	Addr         string
	Serializer   core.Encode
	Deserializer core.Decode
	RetryAttemp  int
}

type ClientProxy interface {
	Send(event string, data interface{}) error
	Broadcast(data interface{}) error
}
