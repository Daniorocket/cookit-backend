package lib

import (
	"net/smtp"
	"os"
)

func CreateEmail(to, subject, msg string) error {

	toList := []string{to}
	host := os.Getenv("Email_HOST")
	username := os.Getenv("EMAIL_LOGIN")
	password := os.Getenv("EMAIL_PASSWORD")
	port := os.Getenv("Email_PORT")
	body := []byte(msg)
	auth := smtp.PlainAuth("", username, password, host)

	if err := smtp.SendMail(host+":"+port, auth, username, toList, body); err != nil {
		return err
	}
	return nil
}
