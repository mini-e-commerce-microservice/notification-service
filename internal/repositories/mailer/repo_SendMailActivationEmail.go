package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"gopkg.in/gomail.v2"
)

func (r *repository) SendMailActivationEmail(ctx context.Context, input SendMailActivationEmailInput) (err error) {
	tmpl := r.template.Lookup(ActivationMail)
	if tmpl == nil {
		return collection.Err(errors.New("mail template not found"))
	}

	tmplBuffer := new(bytes.Buffer)
	tmplData := map[string]any{
		"RecipientName":  input.RecipientName,
		"OTP":            input.OTP,
		"RecipientEmail": input.RecipientEmail,
	}

	if err = tmpl.Execute(tmplBuffer, tmplData); err != nil {
		return collection.Err(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", r.mailerConf.ListEmailAddress.NoReplyEmailAddress)
	m.SetHeader("To", input.RecipientEmail)
	m.SetHeader("Subject", fmt.Sprintf("Verify Your Email"))
	m.SetBody("text/html", tmplBuffer.String())

	err = r.mail.DialAndSend(m)
	if err != nil {
		return collection.Err(err)
	}
	return
}
