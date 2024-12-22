package kafka

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type ConsumerConfig struct {
	GroupID  string
	Assignor string
	Oldest   bool
}

type Consumer struct {
	instance *Kafka
	Group    string
	config   *sarama.Config
	running  bool
	ready    chan bool
	handler  Handler
}

func (k *Kafka) Consumer(opt ConsumerConfig) *Consumer {
	config := sarama.NewConfig()
	config.Version = k.Version

	switch opt.Assignor {
	case sarama.StickyBalanceStrategyName:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case sarama.RoundRobinBalanceStrategyName:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case sarama.RangeBalanceStrategyName:
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", opt.Assignor)
	}

	if opt.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	return &Consumer{
		instance: k,
		Group:    opt.GroupID,
		config:   config,
		ready:    make(chan bool),
		running:  true,
	}
}

func (c *Consumer) Subscribe(topics []string, handler Handler) {
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(c.instance.Brokers, c.Group, c.config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	c.handler = handler
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, c); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	// sigterm := make(chan os.Signal, 1)
	// signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// for c.running {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Println("terminating: context cancelled")
	// 		c.running = false
	// 	case <-sigterm:
	// 		log.Println("terminating: via signal")
	// 		c.running = false
	// 	}
	// }

	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

type Handler func(msg *sarama.ConsumerMessage)

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			consumer.handler(message)
			session.MarkMessage(message, "")
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
