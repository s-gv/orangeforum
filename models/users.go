// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/s-gv/orangeforum/models/db"
	"golang.org/x/crypto/bcrypt"
	"time"
)

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
			if resetDate.After(time.Now().Add(-48 * time.Hour)) {
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
