package kafka_test

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices/kafka"
)

func TestProducer(t *testing.T) {
	producer := kafka.NewProducer(kafka.Options{
		Brokers:   []string{"127.0.0.1:9092"},
		Producers: 10,
		Version:   sarama.DefaultVersion.String(),
	})

	producer.Publish("sarama", "abc")
}
