package kafka

import (
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Config struct {
	microservices.Config
	Broker BrokerConfig
}

func (c Config) IsZero() bool {
	return c.Config.IsZero() && c.Broker.IsZero()
}

type Connect struct {
	Conn    *Kafka
	Module  core.Module
	config  microservices.Config
	timeout time.Duration
}

// Client usage
func NewClient(opt Config) *Connect {
	instance := NewInstance(opt.Broker)
	connect := &Connect{
		Conn:   instance,
		config: microservices.NewConfig(opt.Config),
	}

	return connect
}

func (c *Connect) Config() microservices.Config {
	return c.config
}

func (c *Connect) Publish(event string, data any, headers ...microservices.Header) error {
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

func (c *Connect) Send(path string, request, response any, headers ...microservices.Header) error {
	log.Println("Kafka not support rpc")
	return nil
}

func convertHeaders(headers []*sarama.RecordHeader) microservices.Header {
	header := microservices.Header{}
	for _, h := range headers {
		header[string(h.Key)] = string(h.Value)
	}
	return header
}
