package kafka_test

import (
	"log"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices/kafka"
)

func Test_Consumer(t *testing.T) {
	instance := kafka.New(kafka.Config{
		Brokers: []string{"127.0.0.1:9092"},
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
	producer := instance.Producer(1)
	producer.Publish(&sarama.ProducerMessage{
		Topic: "sarama",
		Value: sarama.StringEncoder("abc"),
	})

	time.Sleep(1000 * time.Millisecond)

}
