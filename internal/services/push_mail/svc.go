package push_mail

import (
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/mailer"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"sync"
)

type service struct {
	rabbitmqRepository rabbitmq.Repository
	mailRepository     mailer.Repository
	wg                 *sync.WaitGroup
	propagators        propagation.TextMapPropagator
}

type NewServiceOpt struct {
	RabbitmqRepository rabbitmq.Repository
	MailRepository     mailer.Repository
	WG                 *sync.WaitGroup
}

func New(opt NewServiceOpt) *service {
	return &service{
		propagators:        otel.GetTextMapPropagator(),
		rabbitmqRepository: opt.RabbitmqRepository,
		mailRepository:     opt.MailRepository,
		wg:                 opt.WG,
	}
}
