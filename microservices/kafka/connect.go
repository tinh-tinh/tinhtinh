package kafka

import (
	"fmt"
	"reflect"
	"time"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/core"
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
	timeout time.Duration
}

// Client usage
func NewClient(opt Options) microservices.ClientProxy {
	instance := NewInstance(opt.Options)
	connect := &Connect{
		Conn:   instance,
		config: microservices.NewConfig(opt.Config),
	}

	return connect
}

func (c *Connect) Config() microservices.Config {
	return c.config
}

func (c *Connect) Publish(event string, data interface{}, headers ...microservices.Header) error {
	payload, err := microservices.EncodeMessage(c, microservices.Message{
		Event:   event,
		Headers: microservices.AssignHeader(c.config.Header, headers...),
		Data:    data,
	})
	if err != nil {
		c.config.ErrorHandler(err)
		return err
	}
	producer := c.Conn.Producer()
	producer.Publish(&sarama.ProducerMessage{
		Topic: event,
		Value: sarama.StringEncoder(string(payload)),
	})
	return nil
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
	consumer := c.Conn.Consumer(ConsumerConfig{
		GroupID:  c.GroupID,
		Assignor: sarama.RangeBalanceStrategyName,
		Oldest:   true,
	})

	var subscribers []*microservices.SubscribeHandler

	store, ok := c.Module.Ref(microservices.STORE).(*microservices.Store)
	if ok && store != nil {
		subscribers = append(subscribers, store.Subscribers...)
	}
	natsStore, ok := c.Module.Ref(microservices.ToTransport(microservices.KAFKA)).(*microservices.Store)
	if ok && natsStore != nil {
		subscribers = append(subscribers, natsStore.Subscribers...)
	}

	// Topics
	topics := Map(subscribers, func(sub *microservices.SubscribeHandler) string {
		return sub.Name
	})
	// handler
	consumer.Subscribe(topics, func(msg *sarama.ConsumerMessage) {
		c.Handler(msg, subscribers)
	})
}

func (c *Connect) Handler(msg *sarama.ConsumerMessage, subscribers []*microservices.SubscribeHandler) {
	sub, ok := Find(subscribers, func(sub *microservices.SubscribeHandler) bool {
		return sub.Name == msg.Topic
	})
	if !ok {
		return
	}

	message := microservices.DecodeMessage(c, msg.Value)
	if reflect.ValueOf(message).IsZero() {
		sub.Handle(c, microservices.Message{
			Data: msg.Value,
		})
	} else {
		sub.Handle(c, message)
	}
}
