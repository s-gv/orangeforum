// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package utils

import (
	"github.com/s-gv/orangeforum/models"
	"log"
	"net/smtp"
)

func SendMail(to string, sub string, body string) {
	go func(to string, sub string, body string) {
		smtpHost := models.Config(models.SMTPHost)
		from := models.Config(models.DefaultFromMail)
		if from != "" && smtpHost != "" {
			smtpUser := models.Config(models.SMTPUser)
			smtpPass := models.Config(models.SMTPPass)
			auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
			msg := []byte("From: " + models.Config(models.ForumName) + "<" + from + ">\r\n" +
				"To: " + to + "\r\n" +
				"Subject: " + sub + "\r\n" +
				"\r\n" +
				body + "\r\n")
			var err error
			if smtpUser != "" {
				err = smtp.SendMail(models.Config(models.SMTPHost)+":"+models.Config(models.SMTPPort), auth, from, []string{to}, msg)
			} else {
				err = smtp.SendMail(models.Config(models.SMTPHost)+":"+models.Config(models.SMTPPort), nil, from, []string{to}, msg)
			}

			if err != nil {
				log.Printf("[ERROR] Error sending mail: %s\n", err)
			}
		} else {
			log.Printf("[ERROR] SMTP not configured.\n")
		}

	}(to, sub, body)
}
