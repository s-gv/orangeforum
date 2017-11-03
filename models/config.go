// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import "github.com/s-gv/orangeforum/models/db"

const (
	ForumName string = "forum_name"
	HeaderMsg string = "header_msg"
	SignupDisabled string = "signup_disabled"
	GroupCreationDisabled string = "group_creation_disabled"
	ImageUploadEnabled string = "image_upload_enabled"
	AllowGroupSubscription string = "allow_group_subscription"
	AllowTopicSubscription string = "allow_topic_subscription"
	DataDir string = "data_dir"
	BodyAppendage string = "body_appendage"
	DefaultFromMail string = "default_from_mail"
	SMTPHost string = "smtp_host"
	SMTPPort string = "smtp_port"
	SMTPUser string = "smtp_user"
	SMTPPass string = "smtp_pass"
	Version string = "version"
)

func IsMigrationNeeded() bool {
	dbver := db.Version()
	return dbver != ModelVersion
}

func WriteConfig(key string, val string) {
	var oldVal string
	if db.QueryRow(`SELECT val FROM configs WHERE name=?;`, key).Scan(&oldVal) == nil {
		if oldVal != val {
			db.Exec(`UPDATE configs SET val=? WHERE name=?;`, val, key)
		}
	} else {
		db.Exec(`INSERT INTO configs(name, val) values(?, ?);`, key, val)
	}
}

func Config(key string) string {
	row := db.QueryRow(`SELECT val FROM configs WHERE name=?;`, key)
	var val string
	if err := row.Scan(&val); err == nil {
		return val
	}
	return "0"
}

func ConfigAllVals() map[string]interface{} {
	vals := map[string]interface{}{
		"forum_name": Config(ForumName),
		"header_msg": Config(HeaderMsg),
		"signup_disabled": Config(SignupDisabled) == "1",
		"group_creation_disabled": Config(GroupCreationDisabled) == "1",
		"image_upload_enabled": Config(ImageUploadEnabled) == "1",
		"allow_group_subscription": Config(AllowGroupSubscription) == "1",
		"allow_topic_subscription": Config(AllowTopicSubscription) == "1",
		"data_dir": Config(DataDir),
		"body_appendage": Config(BodyAppendage),
		"default_from_mail": Config(DefaultFromMail),
		"smtp_host": Config(SMTPHost),
		"smtp_port": Config(SMTPPort),
		"smtp_user": Config(SMTPUser),
		"smtp_pass": Config(SMTPPass),
	}
	return vals
}
