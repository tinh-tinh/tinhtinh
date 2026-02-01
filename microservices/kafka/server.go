package kafka

import (
	"github.com/IBM/sarama"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type ConsumerConnect struct {
	module         core.Module
	config         microservices.Config
	instance       *Kafka
	consumerConfig ConsumerConfig
}

func NewConsumer(broker BrokerConfig, configs ...microservices.Config) *ConsumerConnect {
	config := microservices.ParseConfig(configs...)
	return &ConsumerConnect{
		config:   config,
		instance: NewInstance(broker),
	}
}

func (c *ConsumerConnect) ApplyGroupID(groupID string) *ConsumerConnect {
	c.consumerConfig.GroupID = groupID
	return c
}

func (c *ConsumerConnect) ApplyAssignor(assignor string) *ConsumerConnect {
	c.consumerConfig.Assignor = assignor
	return c
}

func (c *ConsumerConnect) ApplyOldest(oldest bool) *ConsumerConnect {
	c.consumerConfig.Oldest = oldest
	return c
}

func (c *ConsumerConnect) Create(module core.Module) {
	c.module = module
}

func (c *ConsumerConnect) Config() microservices.Config {
	return c.config
}

func (c *ConsumerConnect) Listen() {
	consumer := c.instance.Consumer(c.consumerConfig)

	var subscribers []*microservices.SubscribeHandler

	store, ok := c.module.Ref(microservices.STORE).(*microservices.Store)
	if ok && store != nil {
		subscribers = append(subscribers, store.Subscribers...)
	}
	kafkaStore, ok := c.module.Ref(microservices.ToTransport(microservices.KAFKA)).(*microservices.Store)
	if ok && kafkaStore != nil {
		subscribers = append(subscribers, kafkaStore.Subscribers...)
	}

	topics := common.Map(subscribers, func(sub *microservices.SubscribeHandler) string {
		return sub.Name
	})

	consumer.Subscribe(topics, func(msg *sarama.ConsumerMessage) error {
		return c.Handler(msg, subscribers)
	})
}

func (c *ConsumerConnect) Handler(msg *sarama.ConsumerMessage, subscribers []*microservices.SubscribeHandler) error {
	sub, found := common.Find(subscribers, func(sub *microservices.SubscribeHandler) bool {
		return sub.Name == msg.Topic
	})
	if !found {
		return nil
	}

	message := microservices.DecodeMessage(c, msg.Value)
	if message.IsZero() {
		return sub.Handle(c, microservices.Message{
			Event:   msg.Topic,
			Headers: convertHeaders(msg.Headers),
			Data:    msg.Value,
		})
	} else {
		return sub.Handle(c, message)
	}
}
