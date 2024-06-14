package pkg

import (
	"fmt"
	"gses4_project/internal/models"
	"log"
	"net/smtp"
	"os"
	"sync"
)

type SmtpEmailSender struct {
	from     string
	user     string
	password string
	host     string
	port     string
	auth     smtp.Auth
}

func NewSmtpEmailSender() *SmtpEmailSender {
	smtpSender := &SmtpEmailSender{
		from:     os.Getenv("MAIL_FROM"),
		user:     os.Getenv("MAIL_USER"),
		password: os.Getenv("MAIL_PASSWORD"),
		host:     os.Getenv("MAIL_HOST"),
		port:     os.Getenv("MAIL_PORT"),
	}
	smtpSender.auth = smtp.PlainAuth("", smtpSender.user, smtpSender.password, smtpSender.host)
	return smtpSender
}

func (s *SmtpEmailSender) SendInforming(subscriptions []models.Email, rate float64) {
	subject := "Subject: Daily Exchange Rate\n"
	fromHeader := fmt.Sprintf("From: %s\n", s.from)
	body := fmt.Sprintf("The current exchange rate is: %f UAH/USD", rate)

	var wg sync.WaitGroup

	for _, v := range subscriptions {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			toHeader := fmt.Sprintf("To: %s\n", email)
			msg := []byte(subject + fromHeader + toHeader + "\n" + body)
			err := smtp.SendMail(s.host+":"+s.port, s.auth, s.from, []string{email}, msg)
			if err != nil {
				log.Println(err.Error())
			}
		}(v.Email)
	}

	wg.Wait()
}
