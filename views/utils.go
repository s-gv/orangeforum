// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/smtp"
	"strconv"

	"github.com/golang/glog"
)

func sendMail(from string, to string, sub string, body string, forumName string, smtpHost string, smtpPort int, smtpUser string, smtpPass string) {
	go func(from string, to string, sub string, body string, forumName string, smtpHost string, smtpPort int, smtpUser string, smtpPass string) {
		if from != "" && smtpHost != "" {
			auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
			msg := []byte("From: " + forumName + "<" + from + ">\r\n" +
				"To: " + to + "\r\n" +
				"Subject: " + sub + "\r\n" +
				"\r\n" +
				body + "\r\n")
			var err error
			if smtpUser != "" {
				err = smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, from, []string{to}, msg)
			} else {
				err = smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), nil, from, []string{to}, msg)
			}

			if err != nil {
				glog.Errorf("Error sending mail: %s\n", err)
			}
		} else {
			glog.Infof("[ERROR] SMTP not configured.\n")
		}

	}(from, to, sub, body, forumName, smtpHost, smtpPort, smtpUser, smtpPass)
}
