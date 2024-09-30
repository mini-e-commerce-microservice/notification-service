package push_mail

import (
	"context"
	"errors"
	erabbitmq "github.com/SyaibanAhmadRamadhan/event-bus/rabbitmq"
	"github.com/mini-e-commerce-microservice/notification-service/generated/proto/notification_proto"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/mailer"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/notification-service/internal/util/tracer"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"google.golang.org/protobuf/proto"
)

func (s *service) RunConsumerBackground(ctx context.Context, input RunConsumerBackgroundInput) (err error) {
	consumerOutput, err := s.rabbitmqRepository.Consume(ctx, rabbitmq.ConsumeInput{
		QueueName:    input.QueueName,
		ConsumerName: input.ConsumerName,
		AutoAck:      true,
	})
	if err != nil {
		return tracer.Error(err)
	}

	s.wg.Add(1)

	for {
		select {
		case <-ctx.Done():
			select {
			case d := <-consumerOutput.Messages:
				s.processMessage(d)
			default:
				break
			}
			s.wg.Done()
			return
		case d := <-consumerOutput.Messages:
			s.processMessage(d)
		}
	}
}

func (s *service) processMessage(d amqp.Delivery) {
	carrier := erabbitmq.NewDeliveryMessageCarrier(&d)
	parentCtx := s.propagators.Extract(context.Background(), carrier)

	_, span := otel.Tracer("process push mail").Start(parentCtx, "process push mail notif")
	defer span.End()

	payload := &notification_proto.Notification{}
	err := proto.Unmarshal(d.Body, payload)
	if err != nil {
		tracer.RecordErrorOtel(span, err)
		err = nil
		d.Acknowledger.Ack(d.DeliveryTag, false)
		return
	}

	switch payload.Type {
	case notification_proto.NotificationType_ACTIVATION_EMAIL:
		span.SetAttributes(attribute.String("notification.type", "ACTIVATION_EMAIL"))
		activationEmail, ok := payload.Data.(*notification_proto.Notification_ActivationEmail)
		if !ok {
			tracer.RecordErrorOtel(span, tracer.Error(errors.New("failed type assertion to Notification_ActivationEmail"+
				"but payload type is NotificationType_ACTIVATION_EMAIL")))
			d.Acknowledger.Ack(d.DeliveryTag, false)
			return
		}

		err = s.pushNotifActivationEmail(context.Background(), activationEmail)
		if err != nil {
			tracer.RecordErrorOtel(span, err)
			err = nil
			d.Acknowledger.Nack(d.DeliveryTag, false, true)
			return
		}
	}

	span.SetStatus(codes.Ok, "Push Notification email successfully")
	d.Acknowledger.Ack(d.DeliveryTag, false)
}

func (s *service) pushNotifActivationEmail(ctx context.Context, data *notification_proto.Notification_ActivationEmail) (err error) {
	err = s.mailRepository.SendMailActivationEmail(ctx, mailer.SendMailActivationEmailInput{
		RecipientEmail: data.ActivationEmail.RecipientEmail,
		RecipientName:  data.ActivationEmail.RecipientName,
		OTP:            data.ActivationEmail.OtpCode,
		ExpiredAt:      data.ActivationEmail.ExpiredAt.AsTime(),
	})
	if err != nil {
		return tracer.Error(err)
	}

	return
}
