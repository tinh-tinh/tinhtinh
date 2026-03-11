package kafka

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

type BrokerConfig struct {
	Brokers []string
	Version string
	Topics  []string
	Verbose bool
}

type Kafka struct {
	Brokers []string
	Version sarama.KafkaVersion
}

func NewInstance(config BrokerConfig) *Kafka {
	if config.Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if config.Version == "" {
		config.Version = sarama.DefaultVersion.String()
	}

	version, err := sarama.ParseKafkaVersion(config.Version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	// go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.LstdFlags))
	return &Kafka{
		Brokers: config.Brokers,
		Version: version,
	}
}

func DefaultConfig() BrokerConfig {
	return BrokerConfig{
		Brokers: []string{"localhost:9092"},
		Version: sarama.DefaultVersion.String(),
	}
}

func (c BrokerConfig) IsZero() bool {
	return len(c.Brokers) == 0 && c.Version == "" && len(c.Topics) == 0 && !c.Verbose
}
