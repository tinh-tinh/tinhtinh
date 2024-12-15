package kafka_test

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices/kafka"
)

func Test_Consumer(t *testing.T) {
	consumer := kafka.NewConsumer(kafka.ConsumerOptions{
		Brokers:  []string{"127.0.0.1:9092"},
		Group:    "example",
		Assignor: "range",
		Version:  sarama.DefaultVersion.String(),
	})

	go consumer.Subscribe([]string{"sarama"})
}
