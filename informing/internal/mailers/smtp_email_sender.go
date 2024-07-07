package mailers

import (
	"fmt"
	"net/smtp"
	"os"
)

type SMTPEmailSender struct {
	from     string
	user     string
	password string
	host     string
	port     string
	auth     smtp.Auth
}

func NewSMTPEmailSender() *SMTPEmailSender {
	smtpSender := &SMTPEmailSender{
		from:     os.Getenv("MAIL_FROM"),
		user:     os.Getenv("MAIL_USER"),
		password: os.Getenv("MAIL_PASSWORD"),
		host:     os.Getenv("MAIL_HOST"),
		port:     os.Getenv("MAIL_PORT"),
	}
	smtpSender.auth = smtp.PlainAuth("", smtpSender.user, smtpSender.password, smtpSender.host)
	return smtpSender
}

func (s *SMTPEmailSender) SendInforming(email string, rate float64) error {
	subject := "Subject: Daily Exchange Rate\n"
	fromHeader := fmt.Sprintf("From: %s\n", s.from)
	body := fmt.Sprintf("The current exchange rate is: %f UAH/USD", rate)

	toHeader := fmt.Sprintf("To: %s\n", email)
	msg := []byte(subject + fromHeader + toHeader + "\n" + body)
	err := smtp.SendMail(s.host+":"+s.port, s.auth, s.from, []string{email}, msg)
	if err != nil {
		return err
	}

	return err
}
