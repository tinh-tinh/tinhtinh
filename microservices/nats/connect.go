package nats

import (
	"context"
	"fmt"
	"reflect"
	"time"

	nats_connect "github.com/nats-io/nats.go"
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
	timeout time.Duration
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

func (c *Connect) Timeout(duration time.Duration) microservices.ClientProxy {
	c.timeout = duration
	return c
}

func (c *Connect) Send(event string, data interface{}, headers ...microservices.Header) error {
	payload, err := microservices.EncodeMessage(c, microservices.Message{
		Type:    microservices.RPC,
		Headers: microservices.AssignMap(c.config.Header, headers...),
		Event:   event,
		Data:    data,
	})
	if err != nil {
		c.Serializer(err)
		return err
	}

	ctx, cancel := context.WithTimeout(c.Context, c.timeout)
	defer cancel()

	err = c.emit(ctx, event, payload)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	fmt.Printf("Send payload: %v to event: %s\n", data, event)
	return nil
}

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	payload, err := microservices.EncodeMessage(c, microservices.Message{
		Type:    microservices.PubSub,
		Headers: microservices.AssignMap(c.config.Header, headers...),
		Event:   event,
		Data:    data,
	})
	if err != nil {
		c.Serializer(err)
		return err
	}

	ctx, cancel := context.WithTimeout(c.Context, c.timeout)
	defer cancel()

	err = c.emit(ctx, event, payload)
	if err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	fmt.Printf("Publish payload: %v to event: %s\n", data, event)
	return nil
}

func (c *Connect) emit(ctx context.Context, event string, payload []byte) error {
	done := make(chan error, 1)
	go func() {
		err := c.Conn.Publish(event, payload)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("publish operation timed out: %w", ctx.Err())
	case err := <-done:
		return err
	}
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
				message := microservices.DecodeMessage(c, msg.Data)
				sub.Handle(c, message)
			})
		}
	}

	if store.Subscribers[string(microservices.PubSub)] != nil {
		for _, sub := range store.Subscribers[string(microservices.PubSub)] {
			go c.Conn.Subscribe(sub.Name, func(msg *nats_connect.Msg) {
				message := microservices.DecodeMessage(c, msg.Data)
				sub.Handle(c, message)
			})
		}
	}
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
