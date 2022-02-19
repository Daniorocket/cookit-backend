package lib

import (
	"crypto/tls"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

func CreateEmail(to, subject, msg string) error {

	host := os.Getenv("Email_HOST")
	senderEmail := os.Getenv("EMAIL_LOGIN")
	password := os.Getenv("EMAIL_PASSWORD")
	port, err := strconv.Atoi(os.Getenv("Email_PORT"))
	if err != nil {
		return err
	}
	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", msg)
	d := gomail.NewDialer(host, port, senderEmail, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
