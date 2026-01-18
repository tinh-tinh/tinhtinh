package kafka_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices/kafka"
)

func Test_Consumer(t *testing.T) {
	instance := kafka.NewInstance(kafka.Config{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
	})
	consumer := instance.Consumer(kafka.ConsumerConfig{
		GroupID:  "example",
		Assignor: sarama.RangeBalanceStrategyName,
		Oldest:   true,
	})
	go consumer.Subscribe([]string{"sarama"}, func(msg *sarama.ConsumerMessage) {
		log.Printf("Receive message %v\n", msg)
	})

	time.Sleep(1000 * time.Millisecond)
	producer := instance.Producer()
	producer.Publish(&sarama.ProducerMessage{
		Topic: "sarama",
		Value: sarama.StringEncoder("abc"),
	})

	time.Sleep(1000 * time.Millisecond)
}
