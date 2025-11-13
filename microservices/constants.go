package microservices

import "time"

const (
	RMQ   = "RMQ"
	KAFKA = "KAFKA"
	MQTT  = "MQTT"
	NATS  = "NATS"
	REDIS = "REDIS"
	TCP   = "TCP"
)

const DEFAULT_TIMEOUT = 5 * time.Minute
