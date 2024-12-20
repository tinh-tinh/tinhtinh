package nats

import (
	"context"
	"encoding/json"
	"fmt"

	nats_connect "github.com/nats-io/nats.go"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Connect struct {
	Conn         *nats_connect.Conn
	Module       core.Module
	Context      context.Context
	serializer   core.Encode
	deserializer core.Decode
}

// Client usage
func NewClient(opt microservices.ConnectOptions) microservices.ClientProxy {
	nc, err := nats_connect.Connect(opt.Addr)
	if err != nil {
		panic(err)
	}

	connect := &Connect{
		Conn:         nc,
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
	}
	if opt.Deserializer != nil {
		connect.deserializer = opt.Deserializer
	}
	if opt.Serializer != nil {
		connect.serializer = opt.Serializer
	}

	return connect
}

func (c *Connect) Send(event string, data interface{}) error {
	payload, err := c.serializer(data)
	if err != nil {
		return err
	}

	fmt.Printf("Send payload: %v to event: %s\n", data, event)
	err = c.Conn.Publish(event, payload)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	return nil
}

func (c *Connect) Broadcast(data interface{}) error {
	return c.Send("*", data)
}

// Server usage
func New(module core.ModuleParam, opts ...microservices.ConnectOptions) microservices.Service {
	svc := &Connect{
		Module:       module(),
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
		Context:      context.Background(),
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			svc.serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			svc.deserializer = opts[0].Deserializer
		}

		if opts[0].Addr != "" {
			nc, err := nats_connect.Connect(opts[0].Addr)
			if err != nil {
				panic(err)
			}
			svc.Conn = nc
		}
	}

	return svc
}

func Open(opts ...microservices.ConnectOptions) core.Service {
	svc := &Connect{
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
		Context:      context.Background(),
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			svc.serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			svc.deserializer = opts[0].Deserializer
		}

		if opts[0].Addr != "" {
			nc, err := nats_connect.Connect(opts[0].Addr)
			if err != nil {
				panic(err)
			}
			svc.Conn = nc
		}
	}

	return svc
}

func (c *Connect) Create(module core.Module) {
	c.Module = module
}

func (c *Connect) Listen() {
	fmt.Println("Listening to NATS")
	for _, prd := range c.Module.GetDataProviders() {
		c.Conn.Subscribe(string(prd.GetName()), func(msg *nats_connect.Msg) {
			fmt.Printf("Received message: %s on event: %s\n", string(msg.Data), string(prd.GetName()))
			data := microservices.ParseCtx(string(msg.Data), c)
			prd.GetFactory()(data)
		})
	}
}

func (svc *Connect) Serializer(v interface{}) ([]byte, error) {
	return svc.serializer(v)
}

func (svc *Connect) Deserializer(data []byte, v interface{}) error {
	return svc.deserializer(data, v)
}
