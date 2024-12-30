package kafka

import (
	"fmt"
	"log"
	"reflect"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Options struct {
	microservices.Config
	Options Config
	GroupID string
}

type Connect struct {
	Conn    *Kafka
	Module  core.Module
	config  microservices.Config
	GroupID string
}

// Client usage
func NewClient(opt Options) microservices.ClientProxy {
	instance := NewInstance(opt.Options)
	connect := &Connect{
		Conn:   instance,
		config: opt.Config,
	}

	if reflect.ValueOf(connect.config).IsZero() {
		connect.config = microservices.DefaultConfig()
	}

	return connect
}

func (c *Connect) Serializer(v interface{}) ([]byte, error) {
	return c.config.Serializer(v)
}

func (c *Connect) Deserializer(data []byte, v interface{}) error {
	return c.config.Deserializer(data, v)
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
		fmt.Println(err)
		return err
	}
	fmt.Printf("Send payload: %v to event: %s\n", string(payload), event)
	producer := c.Conn.Producer()
	producer.Publish(&sarama.ProducerMessage{
		Topic: event,
		Value: sarama.StringEncoder(string(payload)),
	})
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
		fmt.Println(err)
		return err
	}
	fmt.Printf("Send payload: %v to event: %s\n", data, event)
	producer := c.Conn.Producer()
	producer.Publish(&sarama.ProducerMessage{
		Topic: event,
		Value: sarama.StringEncoder(string(payload)),
	})
	return nil
}

// Server usage
func New(module core.ModuleParam, opts ...Options) core.Service {
	connect := &Connect{
		config: microservices.DefaultConfig(),
		Module: module(),
	}

	if len(opts) > 0 {
		if opts[0].GroupID != "" {
			connect.GroupID = opts[0].GroupID
		}
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			connect.config = microservices.ParseConfig(opts[0].Config)
		}
		if !reflect.ValueOf(opts[0].Options).IsZero() {
			conn := NewInstance(opts[0].Options)
			connect.Conn = conn
		}
	}

	return connect
}

func Open(opts ...Options) core.Service {
	connect := &Connect{
		config: microservices.DefaultConfig(),
	}

	if len(opts) > 0 {
		if opts[0].GroupID != "" {
			connect.GroupID = opts[0].GroupID
		}
		if !reflect.ValueOf(opts[0].Config).IsZero() {
			connect.config = microservices.ParseConfig(opts[0].Config)
		}
		if !reflect.ValueOf(opts[0].Options).IsZero() {
			conn := NewInstance(opts[0].Options)
			connect.Conn = conn
		}
	}

	return connect
}

func (c *Connect) Create(module core.Module) {
	c.Module = module
}

func (c *Connect) Listen() {
	fmt.Println("Listening to Kafka")
	store := c.Module.Ref(microservices.STORE).(*microservices.Store)
	if store == nil {
		panic("store not found")
	}

	consumer := c.Conn.Consumer(ConsumerConfig{
		GroupID:  c.GroupID,
		Assignor: sarama.RangeBalanceStrategyName,
		Oldest:   true,
	})

	if store.Subscribers[string(microservices.RPC)] != nil {
		for _, sub := range store.Subscribers[string(microservices.RPC)] {
			consumer.Subscribe([]string{sub.Name}, func(msg *sarama.ConsumerMessage) {
				c.Handler(msg, sub)
			})
		}
	}

	if store.Subscribers[string(microservices.PubSub)] != nil {
		for _, sub := range store.Subscribers[string(microservices.PubSub)] {
			consumer.Subscribe([]string{sub.Name}, func(msg *sarama.ConsumerMessage) {
				c.Handler(msg, sub)
			})
		}
	}
}

func (c *Connect) Handler(msg *sarama.ConsumerMessage, sub microservices.SubscribeHandler) {
	fmt.Println(string(msg.Value))
	var message microservices.Message
	err := c.Deserializer(msg.Value, &message)
	if err != nil {
		fmt.Println("Error deserializing message: ", err)
		return
	}

	fmt.Println(message)
	if reflect.ValueOf(message).IsZero() {
		sub.Handle(c, microservices.Message{
			Data: msg.Value,
		})
	} else {
		sub.Handle(c, message)
	}
}

func (c *Connect) ErrorHandler(err error) {
	log.Printf("Error when running tcp: %v\n", err)
}
