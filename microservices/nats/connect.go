package nats

import (
	"context"
	"fmt"
	"reflect"
	"time"

	nats_connect "github.com/nats-io/nats.go"
	"github.com/tinh-tinh/tinhtinh/v2/common/era"
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
		config: microservices.NewConfig(opt.Config),
	}

	return connect
}

func (c *Connect) Config() microservices.Config {
	return c.config
}

func (c *Connect) Send(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultSend(c)(event, data, headers...)
}

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultPublish(c)(event, data, headers...)
}

func (c *Connect) Timeout(duration time.Duration) microservices.ClientProxy {
	c.timeout = duration
	return c
}

func (c *Connect) Emit(event string, message microservices.Message) error {
	payload, err := microservices.EncodeMessage(c, message)
	if err != nil {
		c.config.ErrorHandler(err)
		return err
	}
	if c.timeout > 0 {
		err = era.TimeoutFunc(c.timeout, func(ctx context.Context) error {
			err := c.Conn.Publish(event, payload)
			return err
		})
	} else {
		err = c.Conn.Publish(event, payload)
	}
	if err != nil {
		c.config.ErrorHandler(err)
		return err
	}
	return nil
}

// Server usage
func New(module core.ModuleParam, opts ...Options) microservices.Service {
	connect := &Connect{
		Module: module(),
		config: microservices.DefaultConfig(),
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
		config: microservices.DefaultConfig(),
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

	if store.GetRPC() != nil {
		for _, sub := range store.GetRPC() {
			go c.Conn.Subscribe(sub.Name, func(msg *nats_connect.Msg) {
				c.Handler(msg, sub)
			})
		}
	}

	if store.GetPubSub() != nil {
		for _, sub := range store.GetPubSub() {
			go c.Conn.Subscribe(sub.Name, func(msg *nats_connect.Msg) {
				c.Handler(msg, sub)
			})
		}
	}
}

func (c *Connect) Handler(msg *nats_connect.Msg, sub *microservices.SubscribeHandler) {
	message := microservices.DecodeMessage(c, msg.Data)
	sub.Handle(c, message)
}
