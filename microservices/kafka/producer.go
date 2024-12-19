package kafka

// SIGUSR1 toggle the pause/resume consumption
import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/rcrowley/go-metrics"

	"github.com/IBM/sarama"
)

type Options struct {
	Brokers     []string
	Version     string
	Verbose     bool
	RecordsRate metrics.Meter
	Producers   int64
}

type Producer struct {
	producers int64
	running   bool
	provider  *ProducerProvider
	Record    metrics.Meter
}

func NewProducer(opt Options) *Producer {
	producer := &Producer{
		running:   true,
		Record:    metrics.GetOrRegisterMeter("records.rate", nil),
		producers: opt.Producers,
	}
	log.Println("Starting a new Sarama producer")

	if opt.Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	version, err := sarama.ParseKafkaVersion(opt.Version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	producer.provider = NewProducerProvider(opt.Brokers, func() *sarama.Config {
		config := sarama.NewConfig()
		config.Version = version
		config.Producer.Idempotent = true
		config.Producer.Return.Errors = false
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
		config.Producer.Transaction.Retry.Backoff = 10
		config.Producer.Transaction.ID = "txn_producer"
		config.Net.MaxOpenRequests = 1
		return config
	})
	go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.LstdFlags))

	return producer
}

func (p *Producer) Publish(topic string, message string) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	for i := 0; i < int(p.producers); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				if !p.running {
					cancel()
					return
				}
				select {
				case <-ctx.Done():
					return
				default:
					p.produceTestRecord(topic, message)
				}
			}
		}()
	}

	p.running = false
	wg.Wait()
	cancel()

	p.provider.clear()
}

func (p *Producer) produceTestRecord(topic string, value string) {
	producer := p.provider.borrow()
	defer p.provider.release(producer)

	// Start kafka transaction
	err := producer.BeginTxn()
	if err != nil {
		log.Printf("unable to start txn %s\n", err)
		return
	}

	// Produce some records in transaction
	producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(value)}

	// commit transaction
	err = producer.CommitTxn()
	if err != nil {
		log.Printf("Producer: unable to commit txn %s\n", err)
		for {
			if producer.TxnStatus()&sarama.ProducerTxnFlagFatalError != 0 {
				// fatal error. need to recreate producer.
				log.Printf("Producer: producer is in a fatal state, need to recreate it")
				break
			}
			// If producer is in abortable state, try to abort current transaction.
			if producer.TxnStatus()&sarama.ProducerTxnFlagAbortableError != 0 {
				err = producer.AbortTxn()
				if err != nil {
					// If an error occured just retry it.
					log.Printf("Producer: unable to abort transaction: %+v", err)
					continue
				}
				break
			}
			// if not you can retry
			err = producer.CommitTxn()
			if err != nil {
				log.Printf("Producer: unable to commit txn %s\n", err)
				continue
			}
		}
		return
	}
	p.Record.Mark(p.producers)
}

// pool of producers that ensure transactional-id is unique.
type ProducerProvider struct {
	transactionIdGenerator int32

	producersLock sync.Mutex
	producers     []sarama.AsyncProducer

	producerProvider func() sarama.AsyncProducer
}

func NewProducerProvider(brokers []string, producerConfigurationProvider func() *sarama.Config) *ProducerProvider {
	provider := &ProducerProvider{}
	provider.producerProvider = func() sarama.AsyncProducer {
		config := producerConfigurationProvider()
		suffix := provider.transactionIdGenerator
		// Append transactionIdGenerator to current config.Producer.Transaction.ID to ensure transaction-id uniqueness.
		if config.Producer.Transaction.ID != "" {
			provider.transactionIdGenerator++
			config.Producer.Transaction.ID = config.Producer.Transaction.ID + "-" + fmt.Sprint(suffix)
		}
		producer, err := sarama.NewAsyncProducer(brokers, config)
		if err != nil {
			return nil
		}
		return producer
	}
	return provider
}

func (p *ProducerProvider) borrow() (producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	if len(p.producers) == 0 {
		for {
			producer = p.producerProvider()
			if producer != nil {
				return
			}
		}
	}

	index := len(p.producers) - 1
	producer = p.producers[index]
	p.producers = p.producers[:index]
	return
}

func (p *ProducerProvider) release(producer sarama.AsyncProducer) {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	// If released producer is erroneous close it and don't return it to the producer pool.
	if producer.TxnStatus()&sarama.ProducerTxnFlagInError != 0 {
		// Try to close it
		_ = producer.Close()
		return
	}
	p.producers = append(p.producers, producer)
}

func (p *ProducerProvider) clear() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		producer.Close()
	}
	p.producers = p.producers[:0]
}
