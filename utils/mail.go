package utils

import (
	"net/smtp"
	"github.com/s-gv/orangeforum/models"
	"log"
)

func SendMail(to string, sub string, body string) {
	go func(to string, sub string, body string) {
		smtpHost := models.Config(models.SMTPHost)
		from := models.Config(models.DefaultFromMail)
		if from != "" && smtpHost != "" {
			auth := smtp.PlainAuth("", models.Config(models.SMTPUser), models.Config(models.SMTPPass), smtpHost)

			msg := []byte("From: "+models.Config(models.ForumName)+"<"+from+">\r\n" +
				"To: "+to+"\r\n" +
				"Subject: "+sub+"\r\n" +
				"\r\n" +
				body+"\r\n")
			err := smtp.SendMail(models.Config(models.SMTPHost)+":"+models.Config(models.SMTPPort), auth, from, []string{to}, msg)
			if err != nil {
				log.Printf("[ERROR] Error sending mail: %s\n", err)
			}
		} else {
			log.Printf("[ERROR] SMTP not configured.\n")
		}

	}(to, sub, body)
}