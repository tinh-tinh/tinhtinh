package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/microservices"
)

type Connect struct {
	Conn         *Kafka
	Module       core.Module
	serializer   core.Encode
	deserializer core.Decode
	GroupID      string
}

// Client usage
func NewClient(opt microservices.ConnectOptions) microservices.ClientProxy {
	instance := New(Config{
		Brokers: []string{opt.Addr},
	})
	connect := &Connect{
		Conn:         instance,
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

func (c *Connect) Serializer(v interface{}) ([]byte, error) {
	return c.serializer(v)
}

func (c *Connect) Deserializer(data []byte, v interface{}) error {
	return c.deserializer(data, v)
}

func (c *Connect) Send(event string, data interface{}) error {
	payload, err := c.Serializer(data)
	if err != nil {
		return err
	}
	producer := c.Conn.Producer(1)
	producer.Publish(&sarama.ProducerMessage{
		Topic: event,
		Value: sarama.StringEncoder(string(payload)),
	})
	return nil
}

func (client *Connect) Broadcast(data interface{}) error {
	return client.Send("*", data)
}

// Server usage
func Open(groupID string, opts ...microservices.ConnectOptions) core.Service {
	connect := &Connect{
		serializer:   json.Marshal,
		deserializer: json.Unmarshal,
		GroupID:      groupID,
	}

	if len(opts) > 0 {
		if opts[0].Serializer != nil {
			connect.serializer = opts[0].Serializer
		}

		if opts[0].Deserializer != nil {
			connect.deserializer = opts[0].Deserializer
		}

		if opts[0].Addr != "" {
			conn := New(Config{
				Brokers: []string{opts[0].Addr},
			})
			connect.Conn = conn
		}
	}

	return connect
}

func (c *Connect) Create(module core.Module) {
	c.Module = module
}

func (c *Connect) Listen() {
	consumer := c.Conn.Consumer(ConsumerConfig{
		GroupID:  c.GroupID,
		Assignor: sarama.RangeBalanceStrategyName,
		Oldest:   true,
	})
	for _, prd := range c.Module.GetDataProviders() {
		go consumer.Subscribe([]string{string(prd.GetName())}, func(msg *sarama.ConsumerMessage) {
			data := microservices.ParseCtx(msg.Value, c)
			prd.GetFactory()(data)
		})
	}
}