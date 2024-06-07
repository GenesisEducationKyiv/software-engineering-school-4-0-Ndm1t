package pkg

import (
	"encoding/json"
	"fmt"
	"gses4_project/internal/database"
	"gses4_project/internal/models"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"
)

type GetRateResponse struct {
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func FetchRate() (*float64, error) {
	res, err := http.Get(os.Getenv("API_URL"))
	if err != nil {

		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {

		return nil, err
	}

	var rateList GetRateResponse
	err = json.Unmarshal(body, &rateList)
	if err != nil {
		return nil, err
	}

	rate := rateList.ConversionRates["UAH"]

	return &rate, nil
}

func SendRateEmails() {
	subscriptions, err := ListSubscribed()
	if err != nil {
		log.Println("Failed to list subscriptions:", err)
		return
	}

	rate, err := FetchRate()
	if err != nil {
		log.Println("Failed to fetch rate:", err)
		return
	}

	from := os.Getenv("MAIL_FROM")
	user := os.Getenv("MAIL_USER")
	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	auth := smtp.PlainAuth("", user, password, host)

	subject := "Subject: Daily Exchange Rate\n"
	fromHeader := fmt.Sprintf("From: %s\n", from)
	body := fmt.Sprintf("The current exchange rate is: %f UAH/USD", *rate)

	var wg sync.WaitGroup

	for _, v := range subscriptions {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()
			toHeader := fmt.Sprintf("To: %s\n", email)
			msg := []byte(subject + fromHeader + toHeader + "\n" + body)
			err = smtp.SendMail(host+":"+port, auth, from, []string{email}, msg)
			if err != nil {
				log.Fatalf(err.Error())
			}
		}(v.Email)
	}

	wg.Wait()

}

func ListSubscribed() ([]models.Email, error) {
	var subscriptions []models.Email
	result := database.DB.Find(&subscriptions, "status = ?", models.Subscribed)
	return subscriptions, result.Error
}
