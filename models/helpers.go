// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"github.com/s-gv/orangeforum/models/db"
	"database/sql"
	"time"
	"golang.org/x/crypto/bcrypt"
	"encoding/hex"
	"errors"
)

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

func NumUsers() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM users LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func NumGroups() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM groups LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func NumTopics() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM topics LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func NumComments() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM comments LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func createUser(userName string, passwd string, email string, isSuperAdmin bool) error {
	if passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err == nil {
		r := db.QueryRow(`SELECT username FROM users WHERE username=?;`, userName)
		var tmp string
		if err := r.Scan(&tmp); err == sql.ErrNoRows {
			db.Exec(`INSERT INTO users(username, passwdhash, email, is_superadmin, created_date, updated_date) VALUES(?, ?, ?, ?, ?, ?);`,
				userName, hex.EncodeToString(passwdHash), email, isSuperAdmin, time.Now().Unix(), time.Now().Unix())
		} else {
			return errors.New("Username already exists.")
		}
	} else {
		return err
	}
	return nil
}

func CreateUser(userName string, passwd string, email string) error {
	return createUser(userName, passwd, email, false)
}

func CreateSuperUser(userName string, passwd string) error {
	return createUser(userName, passwd, "", true)
}

func CreateGroupMod(userName string, groupID string) {
	if uid, err := ReadUserIDByName(userName); err == nil {
		db.Exec(`INSERT INTO mods(userid, groupid, created_date) VALUES(?, ?, ?);`, uid, groupID, time.Now().Unix())
	}
}

func CreateGroupAdmin(userName string, groupID string) {
	if uid, err := ReadUserIDByName(userName); err == nil {
		db.Exec(`INSERT INTO admins(userid, groupid, created_date) VALUES(?, ?, ?);`, uid, groupID, time.Now().Unix())
	}
}

func ReadMods(groupID string) []string {
	rows := db.Query(`SELECT users.username FROM users INNER JOIN mods ON users.id=mods.userid WHERE mods.groupid=?;`, groupID)
	var mods []string
	for rows.Next() {
		var mod string
		rows.Scan(&mod)
		mods = append(mods, mod)
	}
	return mods
}


func ReadAdmins(groupID string) []string {
	rows := db.Query(`SELECT users.username FROM users INNER JOIN admins ON users.id=admins.userid WHERE admins.groupid=?;`, groupID)
	var admins []string
	for rows.Next() {
		var admin string
		rows.Scan(&admin)
		admins = append(admins, admin)
	}
	return admins
}

func IsUserGroupAdmin(userID string, groupID string) bool {
	r := db.QueryRow(`SELECT id FROM admins WHERE userid=? AND groupid=?`, userID, groupID)
	var tmp string
	if err := r.Scan(&tmp); err == nil {
		return true
	}
	return false
}

func ReadUserEmail(userName string) string {
	r := db.QueryRow(`SELECT email FROM users WHERE username=?;`, userName)
	var email string
	if err := r.Scan(&email); err == nil {
		return email
	}
	return ""
}

func ReadUserNameByToken(resetToken string) (string, error) {
	if len(resetToken) > 0 {
		r := db.QueryRow(`SELECT username, reset_token_date FROM users WHERE reset_token=?;`, resetToken)
		var userName string
		var rDate int64
		if err := r.Scan(&userName, &rDate); err == nil {
			resetDate := time.Unix(rDate, 0)
			if resetDate.After(time.Now().Add(-48*time.Hour)) {
				return userName, nil
			}
		}
	}
	return "", errors.New("Invalid/Expired reset token.")
}

func ReadUserIDByName(userName string) (int, error) {
	r := db.QueryRow(`SELECT id FROM users WHERE username=?;`, userName)
	var id int
	if err := r.Scan(&id); err == nil {
		return id, nil
	}
	return 0, errors.New("User not found.")
}

func UpdateUserPasswd(userName string, passwd string) error {
	if passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err == nil {
		db.Exec(`UPDATE users SET passwdhash=?, reset_token='', reset_token_date=0 WHERE username=?`, hex.EncodeToString(passwdHash), userName)
	} else {
		return err
	}
	return nil
}

func ProbeUser(userName string) bool {
	r := db.QueryRow(`SELECT username FROM users WHERE username=?;`, userName)
	var tmp string
	if err := r.Scan(&tmp); err == sql.ErrNoRows {
		return false
	}
	return true
}

func ReadGroupIDByName(name string) string {
	r := db.QueryRow(`SELECT id FROM groups WHERE name=?;`, name)
	var id string
	if err := r.Scan(&id); err == nil {
		return id
	}
	return ""
}

