package mail

import "gopkg.in/gomail.v2"

type repository struct {
	mail *gomail.Dialer
}

func New(mail *gomail.Dialer) *repository {
	return &repository{
		mail: mail,
	}
}
