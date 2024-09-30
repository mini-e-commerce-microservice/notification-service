package rabbitmq

import (
	erabbitmq "github.com/SyaibanAhmadRamadhan/event-bus/rabbitmq"
)

type rabbitmq struct {
	r erabbitmq.RabbitMQPubSub
}

func NewRabbitMq(r erabbitmq.RabbitMQPubSub) *rabbitmq {
	return &rabbitmq{
		r: r,
	}
}
