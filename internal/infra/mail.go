package infra

import (
	"github.com/mini-e-commerce-microservice/notification-service/internal/conf"
	"gopkg.in/gomail.v2"
)

func NewMail(cred conf.ConfigMailer) *gomail.Dialer {
	var dialer *gomail.Dialer
	if cred.UsedMailTrap {
		mailTrap := cred.MailTrap
		dialer = gomail.NewDialer(mailTrap.Host, mailTrap.Port, mailTrap.Username, mailTrap.Password)
	} else {
		panic("unknown mailing provider")
	}

	return dialer
}
