// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import "database/sql"

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
	err := DB.Get(&val, "SELECT val FROM configs WHERE name=$1;", key)
	if err != nil {
		return ""
	}
	return val
}

func SetConfigValue(key string, val string) {
	var dummy string
	err := DB.Get(&dummy, "SELECT val FROM configs WHERE name=$1;", key)
	if err == sql.ErrNoRows {
		DB.Exec("INSERT INTO configs(name, val) VALUES($1, $2);", key, val)
	} else {
		DB.Exec("UPDATE configs SET val = $1 WHERE name=$2;", val, key)
	}
}
