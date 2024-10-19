package infra

import (
	"github.com/mini-e-commerce-microservice/notification-service/generated/proto/secret_proto"
	"gopkg.in/gomail.v2"
)

func NewMail(cred *secret_proto.Email) *gomail.Dialer {
	var dialer *gomail.Dialer
	if cred.UseUsedMailTrap {
		mailTrap := cred.MailTrap
		dialer = gomail.NewDialer(mailTrap.Host, int(mailTrap.Port), mailTrap.Username, mailTrap.Password)
	} else {
		panic("unknown mailing provider")
	}

	return dialer
}
