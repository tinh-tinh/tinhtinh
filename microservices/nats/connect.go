package nats

import (
	"context"
	"fmt"
	"reflect"

	nats_connect "github.com/nats-io/nats.go"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Options struct {
	microservices.Config
	Addr string
	nats_connect.Option
}

type Connect struct {
	Conn    *nats_connect.Conn
	Module  core.Module
	Context context.Context
	config  microservices.Config
}

// Client usage
func NewClient(opt Options) microservices.ClientProxy {
	nc, err := nats_connect.Connect(opt.Addr, opt.Option)
	if err != nil {
		panic(err)
	}

	connect := &Connect{
		Conn:   nc,
		config: opt.Config,
	}

	if reflect.ValueOf(connect.config).IsZero() {
		connect.config = microservices.DefaultConfig()
	}

	return connect
}

func (c *Connect) Send(event string, data interface{}, headers ...microservices.Header) error {
	message := microservices.Message{
		Type:    microservices.RPC,
		Headers: common.CloneMap(c.config.Header),
		Event:   event,
		Data:    data,
	}
	if len(headers) > 0 {
		for _, v := range headers {
			common.MergeMaps(message.Headers, v)
		}
	}

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

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	message := microservices.Message{
		Type:    microservices.PubSub,
		Event:   event,
		Data:    data,
		Headers: common.CloneMap(c.config.Header),
	}
	if len(headers) > 0 {
		for _, v := range headers {
			common.MergeMaps(message.Headers, v)
		}
	}

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
func New(module core.ModuleParam, opts ...Options) microservices.Service {
	connect := &Connect{
		Module:  module(),
		config:  microservices.DefaultConfig(),
		Context: context.Background(),
	}

	if len(opts) > 0 {
		if opts[0].Addr != "" {
			nc, err := nats_connect.Connect(opts[0].Addr, opts[0].Option)
			if err != nil {
				panic(err)
			}
			connect.Conn = nc
		}
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			connect.config = microservices.ParseConfig(opts[0].Config)
		}
	}

	return connect
}

func Open(opts ...Options) core.Service {
	connect := &Connect{
		config:  microservices.DefaultConfig(),
		Context: context.Background(),
	}

	if len(opts) > 0 {
		if opts[0].Addr != "" {
			nc, err := nats_connect.Connect(opts[0].Addr, opts[0].Option)
			if err != nil {
				panic(err)
			}
			connect.Conn = nc
		}
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			connect.config = microservices.ParseConfig(opts[0].Config)
		}
	}

	return connect
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

	sub.Handle(c, message)
	// data := microservices.ParseCtx(message.Data, c)
	// factory(data)
}

func (svc *Connect) Serializer(v interface{}) ([]byte, error) {
	return svc.config.Serializer(v)
}

func (svc *Connect) Deserializer(data []byte, v interface{}) error {
	return svc.config.Deserializer(data, v)
}

func (c *Connect) ErrorHandler(err error) {
	c.config.ErrorHandler(err)
}
