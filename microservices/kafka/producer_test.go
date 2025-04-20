package kafka_test

import (
	"os"
	"testing"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices/kafka"
)

func Test_Producer(t *testing.T) {
	instance := kafka.NewInstance(kafka.Config{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
	})
	producer := instance.Producer()
	producer.Publish(&sarama.ProducerMessage{
		Topic: "order.updated",
		Value: sarama.StringEncoder("abc"),
	})
}
