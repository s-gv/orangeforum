// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import "github.com/s-gv/orangeforum/models/db"

const (
	ForumName              string = "forum_name"
	HeaderMsg              string = "header_msg"
	LoginMsg               string = "login_msg"
	SignupMsg              string = "signup_msg"
	CensoredWords          string = "censored_words"
	SignupDisabled         string = "signup_disabled"
	GroupCreationDisabled  string = "group_creation_disabled"
	ImageUploadEnabled     string = "image_upload_enabled"
	AllowGroupSubscription string = "allow_group_subscription"
	AllowTopicSubscription string = "allow_topic_subscription"
	ReadOnlyMode           string = "read_only"
	DataDir                string = "data_dir"
	BodyAppendage          string = "body_appendage"
	DefaultFromMail        string = "default_from_mail"
	SMTPHost               string = "smtp_host"
	SMTPPort               string = "smtp_port"
	SMTPUser               string = "smtp_user"
	SMTPPass               string = "smtp_pass"
	Version                string = "version"
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
	if key == SignupMsg {
		return ""
	}
	if key == LoginMsg {
		return ""
	}
	if key == CensoredWords {
		return ""
	}
	return "0"
}

func ConfigAllVals() map[string]interface{} {
	vals := map[string]interface{}{
		ForumName:              Config(ForumName),
		HeaderMsg:              Config(HeaderMsg),
		LoginMsg:               Config(LoginMsg),
		SignupMsg:              Config(SignupMsg),
		CensoredWords:          Config(CensoredWords),
		SignupDisabled:         Config(SignupDisabled) == "1",
		GroupCreationDisabled:  Config(GroupCreationDisabled) == "1",
		ImageUploadEnabled:     Config(ImageUploadEnabled) == "1",
		AllowGroupSubscription: Config(AllowGroupSubscription) == "1",
		AllowTopicSubscription: Config(AllowTopicSubscription) == "1",
		ReadOnlyMode:           Config(ReadOnlyMode) == "1",
		DataDir:                Config(DataDir),
		BodyAppendage:          Config(BodyAppendage),
		DefaultFromMail:        Config(DefaultFromMail),
		SMTPHost:               Config(SMTPHost),
		SMTPPort:               Config(SMTPPort),
		SMTPUser:               Config(SMTPUser),
		SMTPPass:               Config(SMTPPass),
	}
	return vals
}
