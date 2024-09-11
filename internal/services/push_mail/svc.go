package push_mail

import (
	s3wrapper "github.com/SyaibanAhmadRamadhan/go-s3-wrapper"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/mail"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/rabbitmq"
)

type service struct {
	rabbitmqRepository rabbitmq.Repository
	mailRepository     mail.Repository
	s3                 s3wrapper.S3Client
}

type NewServiceOpt struct {
	RabbitmqRepository rabbitmq.Repository
	MailRepository     mail.Repository
	S3                 s3wrapper.S3Client
}

func New(opt NewServiceOpt) *service {
	return &service{
		rabbitmqRepository: opt.RabbitmqRepository,
		mailRepository:     opt.MailRepository,
		s3:                 opt.S3,
	}
}
