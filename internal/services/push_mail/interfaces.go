package push_mail

import "context"

type Service interface {
	RunConsumerBackground(ctx context.Context, input RunConsumerBackgroundInput) (err error)
}
