package mailers

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/spf13/viper"
	"net/smtp"
)

const (
	emailSendingSuccess = "email_sending_success"
	emailSendingFail    = "email_sending_fail"
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
		from:     viper.GetString("MAIL_FROM"),
		user:     viper.GetString("MAIL_USER"),
		password: viper.GetString("MAIL_PASSWORD"),
		host:     viper.GetString("MAIL_HOST"),
		port:     viper.GetString("MAIL_PORT"),
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
		metrics.GetOrCreateCounter(emailSendingFail).Inc()
		return err
	}
	metrics.GetOrCreateCounter(emailSendingSuccess).Inc()
	return err
}
