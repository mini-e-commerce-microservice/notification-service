package rabbitmq

import (
	"context"
	"github.com/mini-e-commerce-microservice/notification-service/internal/util/tracer"
)

func (r *rabbitmq) Consume(ctx context.Context, input ConsumeInput) (output ConsumeOutput, err error) {
	consume, err := r.ch.ConsumeWithContext(ctx,
		input.QueueName,
		input.ConsumerName,
		input.AutoAck,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return output, tracer.Error(err)
	}

	output = ConsumeOutput{
		Messages: consume,
	}

	return
}
