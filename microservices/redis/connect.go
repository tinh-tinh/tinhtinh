package redis

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	redis_store "github.com/redis/go-redis/v9"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Options struct {
	microservices.Config
	*redis_store.Options
}

type Connect struct {
	Context context.Context
	Module  core.Module
	Conn    *redis_store.Client
	config  microservices.Config
}

// Client usage
func NewClient(opt Options) microservices.ClientProxy {
	conn := redis_store.NewClient(opt.Options)

	if reflect.ValueOf(opt.Config).IsZero() {
		opt.Config = microservices.DefaultConfig()
	} else {
		opt.Config = microservices.ParseConfig(opt.Config)
	}

	connect := &Connect{
		Context: context.Background(),
		Conn:    conn,
		config:  opt.Config,
	}

	return connect
}

func (c *Connect) Send(event string, data interface{}, headers ...microservices.Header) error {
	err := c.emit(event, microservices.Message{
		Type:    microservices.RPC,
		Headers: microservices.AssignHeader(c.config.Header, headers...),
		Event:   event,
		Data:    data,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Send mesage: %v for event: %s\n", data, event)
	return nil
}

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	err := c.emit(event, microservices.Message{
		Type:    microservices.PubSub,
		Headers: microservices.AssignHeader(c.config.Header, headers...),
		Event:   event,
		Data:    data,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Publish mesage: %v for event %s\n", data, event)
	return nil
}

func (c *Connect) emit(event string, message microservices.Message) error {
	payload, err := microservices.EncodeMessage(c, message)
	if err != nil {
		c.Serializer(err)
		return err
	}

	err = c.Conn.Publish(c.Context, event, payload).Err()
	if err != nil {
		return err
	}
	return nil
}

// Server usage
func New(module core.ModuleParam, opts ...Options) microservices.Service {
	connect := &Connect{
		Context: context.Background(),
		Module:  module(),
		config:  microservices.DefaultConfig(),
	}

	if len(opts) > 0 {
		if opts[0].Options != nil {
			conn := redis_store.NewClient(opts[0].Options)
			connect.Conn = conn
		}
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			connect.config = microservices.ParseConfig(opts[0].Config)
		}
	}

	return connect
}

func Open(opts ...Options) core.Service {
	connect := &Connect{
		Context: context.Background(),
		config:  microservices.DefaultConfig(),
	}

	if len(opts) > 0 {
		if opts[0].Options != nil {
			conn := redis_store.NewClient(opts[0].Options)
			connect.Conn = conn
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
	store := c.Module.Ref(microservices.STORE).(*microservices.Store)
	if store == nil {
		panic("store not found")
	}

	if store.Subscribers[string(microservices.RPC)] != nil {
		for _, sub := range store.Subscribers[string(microservices.RPC)] {
			subscriber := c.Conn.Subscribe(c.Context, sub.Name)
			go c.Handler(subscriber, sub)
		}
	}

	if store.Subscribers[string(microservices.PubSub)] != nil {
		var subscriber *redis_store.PubSub
		for _, sub := range store.Subscribers[string(microservices.PubSub)] {
			if strings.HasSuffix(sub.Name, "*") {
				subscriber = c.Conn.PSubscribe(c.Context, sub.Name)
			} else {
				subscriber = c.Conn.Subscribe(c.Context, sub.Name)
			}
			go c.Handler(subscriber, sub)
		}
	}
}

func (c *Connect) Handler(subscriber *redis_store.PubSub, sub microservices.SubscribeHandler) {
	for {
		msg, err := subscriber.ReceiveMessage(c.Context)
		if err != nil {
			return
		}
		message := microservices.DecodeMessage(c, []byte(msg.Payload))
		sub.Handle(c, message)
	}
}

func (c *Connect) Serializer(v interface{}) ([]byte, error) {
	return c.config.Serializer(v)
}

func (c *Connect) Deserializer(data []byte, v interface{}) error {
	return c.config.Deserializer(data, v)
}

func (c *Connect) ErrorHandler(err error) {
	c.config.ErrorHandler(err)
}
