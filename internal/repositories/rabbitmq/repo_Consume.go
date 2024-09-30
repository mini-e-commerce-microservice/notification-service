package rabbitmq

import (
	"context"
	erabbitmq "github.com/SyaibanAhmadRamadhan/event-bus/rabbitmq"
	"github.com/mini-e-commerce-microservice/notification-service/internal/util/tracer"
)

func (r *rabbitmq) Consume(ctx context.Context, input ConsumeInput) (output ConsumeOutput, err error) {
	consume, err := r.r.Subscribe(ctx,
		erabbitmq.SubInput{
			QueueName:    input.QueueName,
			ConsumerName: input.ConsumerName,
			AutoAck:      false,
			Exclusive:    false,
			NoLocal:      false,
			NoWait:       false,
			Args:         nil,
		},
	)
	if err != nil {
		return output, tracer.Error(err)
	}

	output = ConsumeOutput{
		Messages: consume.Deliveries,
	}

	return
}
