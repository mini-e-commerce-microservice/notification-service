package push_mail

import (
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/mailer"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/rabbitmq"
)

type service struct {
	rabbitmqRepository rabbitmq.Repository
	mailRepository     mailer.Repository
}

type NewServiceOpt struct {
	RabbitmqRepository rabbitmq.Repository
	MailRepository     mailer.Repository
}

func New(opt NewServiceOpt) *service {
	return &service{
		rabbitmqRepository: opt.RabbitmqRepository,
		mailRepository:     opt.MailRepository,
	}
}
