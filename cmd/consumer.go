package main

import (
	"context"
	erabbitmq "github.com/SyaibanAhmadRamadhan/event-bus/rabbitmq"
	s3_wrapper_minio "github.com/SyaibanAhmadRamadhan/go-s3-wrapper/minio"
	"github.com/mini-e-commerce-microservice/notification-service/internal/conf"
	"github.com/mini-e-commerce-microservice/notification-service/internal/infra"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/mailer"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/notification-service/internal/services/push_mail"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os/signal"
	"sync"
	"syscall"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "run consumer",
	Run: func(cmd *cobra.Command, args []string) {
		otelConf := conf.LoadOtelConf()
		rabbitMqConf := conf.LoadRabbitMQConf()
		miniConf := conf.LoadMinioConf()
		mailerConf := conf.LoadMailerConf()

		otelClose := infra.NewOtelCollector(otelConf, "notification-service")
		r := erabbitmq.New(rabbitMqConf.Url, erabbitmq.WithOtel(rabbitMqConf.Url))

		mailDialer := infra.NewMail(mailerConf)
		minioClient := infra.NewMinio(miniConf)

		rabbitmqRepository := rabbitmq.NewRabbitMq(r)
		s3 := s3_wrapper_minio.New(minioClient)

		mailRepository := mailer.New(mailer.NewOpt{
			MailerConf:  mailerConf,
			Mail:        mailDialer,
			S3:          s3,
			MinioConfig: miniConf,
		})

		wg := &sync.WaitGroup{}
		pushService := push_mail.New(push_mail.NewServiceOpt{
			RabbitmqRepository: rabbitmqRepository,
			MailRepository:     mailRepository,
			WG:                 wg,
		})

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			err := pushService.RunConsumerBackground(ctx, push_mail.RunConsumerBackgroundInput{
				QueueName:    "notification_type_email",
				ConsumerName: "notification_type_email",
			})
			if err != nil {
				log.Err(err).Msg("consumer notification_type_email background error")
				ctx.Done()
			}
		}()

		<-ctx.Done()
		wg.Wait()

		log.Info().Msg("Received shutdown signal, shutting down server gracefully...")

		if err := otelClose(context.Background()); err != nil {
			log.Err(err).Msg("failed closed otel")
		}

		r.Close()
		log.Info().Msg("Shutdown complete. Exiting.")
		return
	},
}
