package rabbitmq

import "context"

type Repository interface {
	Consume(ctx context.Context, input ConsumeInput) (output ConsumeOutput, err error)
}
