package mailer

import "time"

type SendMailActivationEmailInput struct {
	RecipientEmail string
	RecipientName  string
	OTP            string
	ExpiredAt      time.Time
}
