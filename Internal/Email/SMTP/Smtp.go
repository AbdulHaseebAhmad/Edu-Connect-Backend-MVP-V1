package Smtp

import (
	"fmt"
	"net/smtp"

	Configurator "github.com/AbdulHaseebAhmad/Edu-Connect-Backend-MVP-V1/Internal/Config"
)

type SMTPSender struct {
	Auth smtp.Auth
	Addr string
	From string
}

func NewSMTPSender(cfg *Configurator.Configuration) *SMTPSender {
	auth := smtp.PlainAuth("", cfg.SMTP.Sender, cfg.SMTP.Password, cfg.SMTP.Host)
	addr := cfg.SMTP.Host + ":" + cfg.SMTP.Port
	from := cfg.SMTP.Sender
	return &SMTPSender{
		Auth: auth,
		Addr: addr,
		From: from,
	}
}

func (s *SMTPSender) Send(to, subject, body string) error {
	msg := []byte("Subject: " + subject + "\r\n" + "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" + "\r\n" + body + "\r\n")

	if err := smtp.SendMail(s.Addr, s.Auth, s.From, []string{to}, msg); err != nil {
		return err
	}

	fmt.Println("Email sent to", to)
	return nil
}
