// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/smtp"

	"github.com/golang/glog"
	"github.com/s-gv/orangeforum/models"
)

func sendMail(to string, sub string, body string, forumName string) {
	go func(to string, sub string, body string, forumName string) {
		smtpHost := models.GetConfigValue(models.SMTPHost)
		from := models.GetConfigValue(models.DefaultFromMail)
		if from != "" && smtpHost != "" {
			smtpUser := models.GetConfigValue(models.SMTPUser)
			smtpPass := models.GetConfigValue(models.SMTPPass)
			auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
			msg := []byte("From: " + forumName + "<" + from + ">\r\n" +
				"To: " + to + "\r\n" +
				"Subject: " + sub + "\r\n" +
				"\r\n" +
				body + "\r\n")
			var err error
			if smtpUser != "" {
				err = smtp.SendMail(models.GetConfigValue(models.SMTPHost)+":"+models.GetConfigValue(models.SMTPPort), auth, from, []string{to}, msg)
			} else {
				err = smtp.SendMail(models.GetConfigValue(models.SMTPHost)+":"+models.GetConfigValue(models.SMTPPort), nil, from, []string{to}, msg)
			}

			if err != nil {
				glog.Errorf("Error sending mail: %s\n", err)
			}
		} else {
			glog.Infof("[ERROR] SMTP not configured.\n")
		}

	}(to, sub, body, forumName)
}
