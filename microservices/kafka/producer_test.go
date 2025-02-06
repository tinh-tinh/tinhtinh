package kafka_test

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices/kafka"
)

func Test_Producer(t *testing.T) {
	instance := kafka.NewInstance(kafka.Config{
		Brokers: []string{"127.0.0.1:9092"},
	})
	producer := instance.Producer()
	producer.Publish(&sarama.ProducerMessage{
		Topic: "order.updated",
		Value: sarama.StringEncoder("abc"),
	})
}
