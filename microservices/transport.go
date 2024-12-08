package microservices

type Transport string

const (
	TCP      Transport = "tcp"
	GRPC     Transport = "grpc"
	KAFKA    Transport = "kafka"
	RABBITMQ Transport = "rabbitmq"
	REDIS    Transport = "redis"
)

type Package struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
