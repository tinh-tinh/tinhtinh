package kafka

import (
	"context"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/IBM/sarama"
	"github.com/rcrowley/go-metrics"
)

type ConsumerOptions struct {
	Brokers     []string
	Group       string
	Version     string
	Verbose     bool
	Assignor    string
	RecordsRate metrics.Meter
	Oldest      bool
}
type Consumer struct {
	running bool
	Brokers []string
	Group   string
	config  *sarama.Config
	ready   chan bool
}

func NewConsumer(opt ConsumerOptions) *Consumer {
	consumer := &Consumer{
		running: true,
		Brokers: opt.Brokers,
		ready:   make(chan bool),
		Group:   opt.Group,
	}

	log.Println("Starting a new Sarama consumer")

	if opt.Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion(opt.Version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	config := sarama.NewConfig()
	config.Version = version

	switch opt.Assignor {
	case sarama.StickyBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case sarama.RangeBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	case sarama.RoundRobinBalanceStrategyName:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", opt.Assignor)
	}

	if opt.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	consumer.config = config

	return consumer
}

func (consumer *Consumer) Subscribe(topics []string) {
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(consumer.Brokers, consumer.Group, consumer.config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	consumptionIsPaused := false
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	toggleConsumptionFlow(client, &consumptionIsPaused)

	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		client.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}
