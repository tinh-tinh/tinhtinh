package kafka_test

import (
	"errors"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices/kafka"
)

func Test_Consumer(t *testing.T) {
	instance := kafka.NewInstance(kafka.BrokerConfig{
		Brokers: []string{"127.0.0.1:9092"},
	})
	consumer := instance.Consumer(kafka.ConsumerConfig{
		GroupID:  "example",
		Assignor: sarama.RangeBalanceStrategyName,
		Oldest:   true,
	})
	go consumer.Subscribe([]string{"sarama"}, func(msg *sarama.ConsumerMessage) error {
		log.Printf("Receive message %v\n", msg)
		return nil
	})

	time.Sleep(1000 * time.Millisecond)
	producer := instance.Producer()
	producer.Publish(&sarama.ProducerMessage{
		Topic: "sarama",
		Value: sarama.StringEncoder("abc"),
	})

	time.Sleep(1000 * time.Millisecond)
}

func Test_HandleFailedNotCommit(t *testing.T) {
	instance := kafka.NewInstance(kafka.BrokerConfig{
		Brokers: []string{"127.0.0.1:9092"},
	})

	var count int32 = 0

	consumer := instance.Consumer(kafka.ConsumerConfig{
		GroupID:  "retry-test-group",
		Assignor: sarama.RangeBalanceStrategyName,
		Oldest:   true,
	})

	// Handler that always fails - offset should NOT be committed, message should be retried
	go consumer.Subscribe([]string{"retry-test-topic"}, func(msg *sarama.ConsumerMessage) error {
		atomic.AddInt32(&count, 1)
		currentCount := atomic.LoadInt32(&count)
		log.Printf("Attempt %d: Receive message %s\n", currentCount, string(msg.Value))
		// Always return error - offset should NOT be committed
		return errors.New("simulated failure")
	})

	time.Sleep(1000 * time.Millisecond)

	// Publish a message
	producer := instance.Producer()
	producer.Publish(&sarama.ProducerMessage{
		Topic: "retry-test-topic",
		Value: sarama.StringEncoder("test-retry-message"),
	})

	// Wait for Kafka to retry multiple times
	time.Sleep(5000 * time.Millisecond)

	finalCount := atomic.LoadInt32(&count)
	log.Printf("Total attempts: %d\n", finalCount)

	// If offset is NOT committed on error, Kafka should re-deliver the message
	// So we expect count >= 2 (multiple retry attempts)
	if finalCount < 2 {
		t.Fatalf("Expected at least 2 retry attempts, got %d", finalCount)
	}
}
