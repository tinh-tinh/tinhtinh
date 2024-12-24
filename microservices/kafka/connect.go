package kafka

import (
	"encoding/json"
	"fmt"

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
	message := microservices.Message{Type: microservices.RPC, Event: event, Data: data}
	payload, err := c.Serializer(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Send payload: %v to event: %s\n", string(payload), event)
	producer := c.Conn.Producer(10)
	producer.Publish(&sarama.ProducerMessage{
		Topic: event,
		Value: sarama.StringEncoder(string(payload)),
	})
	return nil
}

func (c *Connect) Publish(event string, data interface{}) error {
	message := microservices.Message{Type: microservices.PubSub, Event: event, Data: data}
	payload, err := c.Serializer(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Send payload: %v to event: %s\n", data, event)
	producer := c.Conn.Producer(10)
	producer.Publish(&sarama.ProducerMessage{
		Topic: event,
		Value: sarama.StringEncoder(string(payload)),
	})
	return nil
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
			go consumer.Subscribe([]string{sub.Name}, func(msg *sarama.ConsumerMessage) {
				c.Handler(msg, sub.Factory)
			})
		}
	}

	if store.Subscribers[string(microservices.PubSub)] != nil {
		for _, sub := range store.Subscribers[string(microservices.PubSub)] {
			go consumer.Subscribe([]string{sub.Name}, func(msg *sarama.ConsumerMessage) {
				c.Handler(msg, sub.Factory)
			})
		}
	}
}

func (c *Connect) Handler(msg *sarama.ConsumerMessage, factory microservices.Factory) {
	fmt.Println(string(msg.Value))
	var message microservices.Message
	err := c.Deserializer(msg.Value, &message)
	if err != nil {
		fmt.Println("Error deserializing message: ", err)
		return
	}

	fmt.Println(message)
	data := microservices.ParseCtx(message.Data, c)
	factory(data)
}
