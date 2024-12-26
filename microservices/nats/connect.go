package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	message := microservices.Message{Type: microservices.RPC, Event: event, Data: data}
	payload, err := c.Serializer(message)
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

func (c *Connect) Publish(event string, data interface{}) error {
	message := microservices.Message{Type: microservices.RPC, Event: event, Data: data}
	payload, err := c.Serializer(message)
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
	store := c.Module.Ref(microservices.STORE).(*microservices.Store)
	if store == nil {
		panic("store not found")
	}

	if store.Subscribers[string(microservices.RPC)] != nil {
		for _, sub := range store.Subscribers[string(microservices.RPC)] {
			go c.Conn.Subscribe(sub.Name, func(msg *nats_connect.Msg) {
				c.Handler(msg, sub)
			})
		}
	}

	if store.Subscribers[string(microservices.PubSub)] != nil {
		for _, sub := range store.Subscribers[string(microservices.PubSub)] {
			go c.Conn.Subscribe(sub.Name, func(msg *nats_connect.Msg) {
				c.Handler(msg, sub)
			})
		}
	}
}

func (c *Connect) Handler(msg *nats_connect.Msg, sub microservices.SubscribeHandler) {
	var message microservices.Message
	err := c.Deserializer([]byte(msg.Data), &message)
	if err != nil {
		fmt.Println("Error deserializing message: ", err)
		return
	}

	sub.Handle(c, message.Data)
	// data := microservices.ParseCtx(message.Data, c)
	// factory(data)
}

func (svc *Connect) Serializer(v interface{}) ([]byte, error) {
	return svc.serializer(v)
}

func (svc *Connect) Deserializer(data []byte, v interface{}) error {
	return svc.deserializer(data, v)
}

func (c *Connect) ErrorHandler(err error) {
	log.Printf("Error when running tcp: %v\n", err)
}
