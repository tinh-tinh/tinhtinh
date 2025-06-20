package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/rcrowley/go-metrics"
	"github.com/tinh-tinh/tinhtinh/v2/common"
)

type Producer struct {
	numProducers int
	txnProducer  *txnProducer
}

func (k *Kafka) Producer(configs ...*sarama.Config) *Producer {
	tProducer := &txnProducer{
		recordsRate: metrics.GetOrRegisterMeter("records.rate", nil),
	}
	tProducer.producerProvider = func() sarama.AsyncProducer {
		var cfg *sarama.Config
		if len(configs) > 0 {
			cfg = common.MergeStruct(configs...)
		} else {
			cfg = defaultProducerConfig(k.Version)
		}
		suffix := tProducer.transactionIdGenerator
		if cfg.Producer.Transaction.ID != "" {
			tProducer.transactionIdGenerator++
			cfg.Producer.Transaction.ID = cfg.Producer.Transaction.ID + "-" + string(suffix)
		}

		producer, err := sarama.NewAsyncProducer(k.Brokers, cfg)
		if err != nil {
			return nil
		}
		return producer
	}
	return &Producer{
		txnProducer:  tProducer,
		numProducers: 1,
	}
}

func (p *Producer) Publish(msg *sarama.ProducerMessage) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var once sync.Once

	for range p.numProducers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				// Context was cancelled before this goroutine started
				return
			default:
				p.txnProducer.producer(msg) // Send the message
				once.Do(func() {
					cancel() // Cancel context after first successful publish
				})
			}
		}()
	}

	wg.Wait()
	p.txnProducer.clear()
}

func defaultProducerConfig(version sarama.KafkaVersion) *sarama.Config {
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
}

type txnProducer struct {
	transactionIdGenerator int32
	producersLock          sync.Mutex
	producers              []sarama.AsyncProducer
	producerProvider       func() sarama.AsyncProducer
	recordsRate            metrics.Meter
}

func (p *txnProducer) producer(msg *sarama.ProducerMessage) {
	producer := p.borrow()
	defer p.release(producer)

	err := producer.BeginTxn()
	if err != nil {
		fmt.Printf("Error beginning transaction: %v\n", err)
		return
	}

	producer.Input() <- msg

	err = producer.CommitTxn()
	if err != nil {
		fmt.Printf("Producer: unable to commit txn %s\n", err)
		for {
			if producer.TxnStatus()&sarama.ProducerTxnFlagFatalError != 0 {
				fmt.Println("Producer: producer is in a fatal state, need to recreate it")
				break
			}
			if producer.TxnStatus()&sarama.ProducerTxnFlagAbortableError != 0 {
				err = producer.AbortTxn()
				if err != nil {
					fmt.Printf("Producer: unable to abort transaction: %+v", err)
					continue
				}
				break
			}
			err = producer.CommitTxn()
			if err != nil {
				fmt.Printf("Producer: unable to commit txn %s\n", err)
				continue
			}
		}
		return
	}
	p.recordsRate.Mark(1)
}

func (p *txnProducer) borrow() (producer sarama.AsyncProducer) {
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

func (p *txnProducer) release(producer sarama.AsyncProducer) {
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

func (p *txnProducer) clear() {
	p.producersLock.Lock()
	defer p.producersLock.Unlock()

	for _, producer := range p.producers {
		producer.Close()
	}
	p.producers = p.producers[:0]
}
