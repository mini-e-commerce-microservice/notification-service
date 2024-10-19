package mailer

import (
	"context"
	s3wrapper "github.com/SyaibanAhmadRamadhan/go-s3-wrapper"
	"github.com/mini-e-commerce-microservice/notification-service/generated/proto/secret_proto"
	"gopkg.in/gomail.v2"
	"html/template"
	"io"
)

type repository struct {
	mail       *gomail.Dialer
	mailerConf *secret_proto.Email
	template   *template.Template
}

type NewOpt struct {
	MailerConf  *secret_proto.Email
	Mail        *gomail.Dialer
	S3          s3wrapper.S3Client
	MinioConfig *secret_proto.Minio
}

func New(opt NewOpt) *repository {
	r := repository{
		mail:       opt.Mail,
		mailerConf: opt.MailerConf,
		template:   new(template.Template),
	}

	templateList := map[string]string{
		ActivationMail: opt.MailerConf.TemplateHtml.ActivationEmailOtp,
	}

	ctx := context.Background()
	for tmplName, tmplPath := range templateList {
		objectOutput, err := opt.S3.GetObject(ctx, s3wrapper.GetObjectInput{
			ObjectName: tmplPath,
			BucketName: opt.MinioConfig.PrivateBucket,
		})
		if err != nil {
			panic(err)
		}

		tmplData, err := io.ReadAll(objectOutput.Object)
		if err != nil {
			panic(err)
		}

		tmpl, err := template.New(tmplName).Parse(string(tmplData))
		if err != nil {
			panic(err)
		}

		r.template = tmpl

		err = objectOutput.Object.Close()
		if err != nil {
			panic(err)
		}
	}

	return &r
}
