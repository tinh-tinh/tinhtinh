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

	if reflect.ValueOf(opt.Config).IsZero() {
		opt.Config = microservices.DefaultConfig()
	}

	connect := &Connect{
		client: conn,
		config: opt.Config,
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

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	token := c.client.Publish(event, 0, false, payload)
	token.Wait()

	c.client.Disconnect(250)
	return nil
}

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	message := microservices.Message{
		Type:    microservices.PubSub,
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
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	token := c.client.Publish(event, 0, false, payload)
	token.Wait()

	c.client.Disconnect(250)
	return nil
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

	if store.Subscribers[string(microservices.RPC)] != nil {
		for _, sub := range store.Subscribers[string(microservices.RPC)] {
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

	if store.Subscribers[string(microservices.PubSub)] != nil {
		for _, sub := range store.Subscribers[string(microservices.PubSub)] {
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

func (c *Connect) handler(msg mqtt_store.Message, sub microservices.SubscribeHandler) {
	var message microservices.Message
	err := c.Deserializer(msg.Payload(), &message)
	if err != nil {
		fmt.Println("Error deserializing message: ", err)
		return
	}

	sub.Handle(c, message)
}
