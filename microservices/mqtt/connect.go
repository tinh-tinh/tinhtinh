package mqtt

import (
	"fmt"
	"reflect"

	mqtt_store "github.com/eclipse/paho.mqtt.golang"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Connect struct {
	Module core.Module
	client mqtt_store.Client
	config microservices.Config
}

type Options struct {
	*mqtt_store.ClientOptions
	microservices.Config
}

func NewClient(opt Options) microservices.ClientProxy {
	conn := mqtt_store.NewClient(opt.ClientOptions)

	connect := &Connect{
		client: conn,
		config: microservices.NewConfig(opt.Config),
	}

	return connect
}

func (c *Connect) Config() microservices.Config {
	return c.config
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

func (c *Connect) Send(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultSend(c)(event, data, headers...)
}

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	return microservices.DefaultPublish(c)(event, data, headers...)
}

func (c *Connect) Emit(event string, message microservices.Message) error {
	payload, err := microservices.EncodeMessage(c, message)
	if err != nil {
		c.ErrorHandler(err)
		return err
	}
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		c.ErrorHandler(err)
		return token.Error()
	}

	token := c.client.Publish(event, 0, false, payload)
	token.Wait()

	c.client.Disconnect(250)
	return nil
}

// Server usage
func New(module core.ModuleParam, opts ...Options) microservices.Service {
	connect := &Connect{
		Module: module(),
		config: microservices.DefaultConfig(),
	}

	if len(opts) > 0 {
		if opts[0].ClientOptions != nil {
			conn := mqtt_store.NewClient(opts[0].ClientOptions)
			connect.client = conn
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
		if opts[0].ClientOptions != nil {
			conn := mqtt_store.NewClient(opts[0].ClientOptions)
			connect.client = conn
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

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if store.GetRPC() != nil {
		for _, sub := range store.GetRPC() {
			token := c.client.Subscribe(sub.Name, 0, func(client mqtt_store.Client, m mqtt_store.Message) {
				c.handler(m, sub)
			})
			token.Wait()
			if token.Error() != nil {
				fmt.Println(token.Error(), common.GetStructName(c.Module))
				continue
			}
		}
	}

	if store.GetPubSub() != nil {
		for _, sub := range store.GetPubSub() {
			token := c.client.Subscribe(sub.Name, 0, func(client mqtt_store.Client, m mqtt_store.Message) {
				c.handler(m, sub)
			})
			token.Wait()
			if token.Error() != nil {
				fmt.Println(token.Error(), common.GetStructName(c.Module))
				continue
			}
		}
	}
}

func (c *Connect) handler(msg mqtt_store.Message, sub *microservices.SubscribeHandler) {
	message := microservices.DecodeMessage(c, msg.Payload())
	sub.Handle(c, message)
}
