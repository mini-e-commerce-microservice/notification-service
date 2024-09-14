package push_mail

import (
	"context"
	"errors"
	"github.com/mini-e-commerce-microservice/notification-service/generated/proto/notification_proto"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/mailer"
	"github.com/mini-e-commerce-microservice/notification-service/internal/repositories/rabbitmq"
	"github.com/mini-e-commerce-microservice/notification-service/internal/util/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

	for d := range consumerOutput.Messages {
		newCtx, span := otel.Tracer("rabbitmq").Start(context.Background(), "consumer message", trace.WithAttributes(
			attribute.String("rabbitmq.correlation_id", d.CorrelationId),
			attribute.String("rabbitmq.exchange", d.Exchange),
			attribute.String("rabbitmq.routing_key", d.Exchange),
			attribute.String("rabbitmq.content_type", d.ContentType),
			attribute.String("rabbitmq.content_encoding", d.ContentEncoding),
			attribute.String("rabbitmq.routing_key", d.RoutingKey),
			attribute.String("rabbitmq.message_id", d.MessageId),
		))

		payload := &notification_proto.Notification{}
		err = proto.Unmarshal(d.Body, payload)
		if err != nil {
			tracer.RecordErrorOtel(span, err)
			err = nil
		}

		switch payload.Type {
		case notification_proto.NotificationType_ACTIVATION_EMAIL:
			err = s.pushNotifActivationEmail(newCtx, payload.Data)
			if err != nil {
				tracer.RecordErrorOtel(span, err)
			}
		}

		span.End()
	}

	return
}

func (s *service) pushNotifActivationEmail(ctx context.Context, data any) (err error) {

	activationEmail, ok := data.(*notification_proto.Notification_ActivationEmail)
	if !ok {
		return tracer.Error(errors.New("failed type assertion to Notification_ActivationEmail" +
			"but payload type is NotificationType_ACTIVATION_EMAIL"))
	}

	err = s.mailRepository.SendMailActivationEmail(ctx, mailer.SendMailActivationEmailInput{
		RecipientEmail: activationEmail.ActivationEmail.RecipientEmail,
		RecipientName:  activationEmail.ActivationEmail.RecipientName,
		OTP:            activationEmail.ActivationEmail.OtpCode,
		ExpiredAt:      activationEmail.ActivationEmail.ExpiredAt.AsTime(),
	})
	if err != nil {
		return tracer.Error(err)
	}

	return
}
