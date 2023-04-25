package config

import (
	"fmt"
	"net/smtp"
	"strconv"
)

func SendMail(configuration Config, toEmail string, subject string, message string) error {
	host := configuration.Get("MAIL_SMTP_HOST")
	port, err := strconv.Atoi(configuration.Get("MAIL_SMTP_PORT"))
	if err != nil {
		return err
	}
	user := configuration.Get("MAIL_AUTH_USER")
	passowrd := configuration.Get("MAIL_AUTH_PASSWORD")

	smtpAuth := smtp.PlainAuth("", user, passowrd, host)
	smtpAddress := fmt.Sprintf("%s:%d", host, port)

	fromEmail := "wids@ganteng.com"
	body := "From: " + fromEmail + "\n" +
		"To: " + toEmail + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	err = smtp.SendMail(smtpAddress, smtpAuth, fromEmail, []string{toEmail}, []byte(body))
	if err != nil {
		return err
	}

	return nil
}
