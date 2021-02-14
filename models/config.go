// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

const (
	ForumName = "forum_name"

	SMTPHost        = "smtp_host"
	SMTPUser        = "smtp_user"
	SMTPPort        = "smtp_port"
	SMTPPass        = "smtp_pass"
	DefaultFromMail = "default_from_mail"
)

func GetConfigValue(key string) string {
	var val string
	err := DB.Get(&val, "SELECT value FROM config WHERE key=$1;", key)
	if err != nil {
		return ""
	}
	return val
}
