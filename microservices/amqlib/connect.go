package amqlib

import (
	"context"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tinh-tinh/tinhtinh/microservices"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

type Options struct {
	microservices.Config
	Addr string
}

type Connect struct {
	Conn    *amqp091.Connection
	Module  core.Module
	Context context.Context
	config  microservices.Config
	timeout time.Duration
}
