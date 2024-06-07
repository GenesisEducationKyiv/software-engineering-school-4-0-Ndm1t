package crons

import (
	"gses4_project/internal/pkg"
	"log"
)

func SendRateEmails() {
	log.Println("Running daily email job")
	pkg.SendRateEmails()
}
