package kafka

import (
	"log"
	"os"

	"github.com/IBM/sarama"
)

type Config struct {
	Brokers []string
	Version string
	Topics  []string
	Verbose bool
}

type Kafka struct {
	Brokers []string
	Version sarama.KafkaVersion
}

func NewInstance(config Config) *Kafka {
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

	if len(config.Brokers) > 0 {
		if config.Brokers[0] == "" {
			config.Brokers[0] = "localhost:9092"
		}
	}

	// go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.LstdFlags))
	return &Kafka{
		Brokers: config.Brokers,
		Version: version,
	}
}

func DefaultConfig() Config {
	return Config{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Version: sarama.DefaultVersion.String(),
	}
}
