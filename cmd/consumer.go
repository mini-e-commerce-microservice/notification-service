package main

import (
	"context"
	s3_wrapper_minio "github.com/SyaibanAhmadRamadhan/go-s3-wrapper/minio"
	"github.com/mini-e-commerce-microservice/notification-service/internal/conf"
	"github.com/mini-e-commerce-microservice/notification-service/internal/infra"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/mail"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/notification-service/internal/services/push_mail"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "run consumer",
	Run: func(cmd *cobra.Command, args []string) {
		conf.Init()

		otelClose := infra.NewOtelCollector(conf.GetConfig().OpenTelemetry)
		_, ch, closeRabbitmq := infra.NewRabbitMq(conf.GetConfig().RabbitMQ)
		mailDialer := infra.NewMail(conf.GetConfig().Mailing)
		minioClient := infra.NewMinio(conf.GetConfig().Minio)

		mailRepository := mail.New(mailDialer)
		rabbitmqRepository := rabbitmq.NewRabbitMq(ch)
		s3 := s3_wrapper_minio.New(minioClient)

		pushService := push_mail.New(push_mail.NewServiceOpt{
			RabbitmqRepository: rabbitmqRepository,
			MailRepository:     mailRepository,
			S3:                 s3,
		})

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			err := pushService.RunConsumerBackground(ctx, push_mail.RunConsumerBackgroundInput{
				QueueName:    "notification_type_email",
				ConsumerName: "notification_type_email",
			})
			if err != nil {
				panic(err)
			}
		}()

		<-ctx.Done()

		log.Info().Msg("Received shutdown signal, shutting down server gracefully...")

		if err := otelClose(context.TODO()); err != nil {
			panic(err)
		}

		if err := closeRabbitmq(context.TODO()); err != nil {
			panic(err)
		}

		log.Info().Msg("Shutdown complete. Exiting.")
		return
	},
}
