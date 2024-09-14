package mailer

import "context"

type Repository interface {
	SendMailActivationEmail(ctx context.Context, input SendMailActivationEmailInput) (err error)
}
