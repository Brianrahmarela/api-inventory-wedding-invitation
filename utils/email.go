package utils

import (
	"net/smtp"
	"os"
	"strings"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_GMAIL_ACCOUNT")

	password := os.Getenv("SMTP_GMAIL_APP_PASSWORD") // App Password dari Gmail (bukan password biasa)

	// SMTP config
	smtpHost := os.Getenv("SMTP_GMAIL_HOST")
	smtpPort := os.Getenv("SMTP_GMAIL_PORT")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}

func FormatGuestLinks(guestLinks []string) string {
	return "Berikut link undangan tamu anda:\n\n" + strings.Join(guestLinks, "\n")
}
