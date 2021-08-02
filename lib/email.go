package lib

import (
	"net/smtp"
	"os"
)

func CreateEmail(to, subject, msg string) error {

	toList := []string{to}
	host := "smtp.gmail.com"
	username := os.Getenv("EMAIL_LOGIN")
	password := os.Getenv("EMAIL_PASSWORD")
	port := "587"
	body := []byte(msg)
	auth := smtp.PlainAuth("", username, password, host)

	if err := smtp.SendMail(host+":"+port, auth, username, toList, body); err != nil {
		return err
	}
	return nil
}
