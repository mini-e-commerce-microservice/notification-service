package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

type ConsumeInput struct {
	QueueName    string
	ConsumerName string
	AutoAck      bool
}

type ConsumeOutput struct {
	Messages <-chan amqp.Delivery
}
