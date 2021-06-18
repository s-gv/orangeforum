// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uuid.UUID `db:"id"`
	Username    string
	PasswdHash  string `db:"passwd_hash"`
	Email       string
	LoggedOutAt time.Time `db:"logout_at"`
}

func createUser(domainName string, email string, userName string, passwd string, isSuperUser bool) error {
	domainID := GetDomainIDByName(domainName)
	if domainID == nil {
		return errors.New("Invalid domain")
	}
	passwdHash := hashPassword(passwd)
	_, err := DB.Exec(`INSERT INTO users(domain_id, email, username, passwd_hash) VALUES($1, $2, $3, $4);`,
		domainID, email, userName, passwdHash,
	)
	return err
}

func CreateUser(domainName string, email string, userName string, passwd string) error {
	return createUser(domainName, email, userName, passwd, false)
}

func CreateSuperUser(domainName string, email string, userName string, passwd string) error {
	return createUser(domainName, email, userName, passwd, true)
}

func ChangePasswd(domainName string, email string, passwd string) error {
	domainID := GetDomainIDByName(domainName)
	if domainID == nil {
		return errors.New("Invalid domain")
	}
	passwdHash := hashPassword(passwd)
	_, err := DB.Exec(`UPDATE users SET passwd_hash = $1 WHERE domain_id = $2 AND email = $3;`,
		passwdHash, domainID, email,
	)
	return err
}

func GetUserByID(id uuid.UUID) *User {
	user := User{}
	err := DB.Get(&user, "SELECT id, username, passwd_hash, email, logout_at FROM users WHERE id=$1;", id)
	if err != nil {
		return nil
	}
	return &user
}

func GetUsersByEmail(email string) *[]User {
	users := []User{}
	DB.Select(&users, "SELECT id, username, passwd_hash, email, logout_at FROM users WHERE email=$1;", email)
	return &users
}

func GetUserByPasswd(username, passwd string) *User {
	user := User{}
	err := DB.Get(&user, "SELECT id, username, passwd_hash, email, logout_at FROM users WHERE username=$1;", username)
	if err != nil {
		return nil
	}
	if !checkPasswordHash(passwd, user.PasswdHash) {
		return nil
	}
	return &user
}

func GetUserByOneTimeToken(oneTimeToken string) *User {
	user := User{}
	err := DB.Get(&user, "SELECT id, username, passwd_hash, email, logout_at FROM users WHERE onetime_login_token=$1;", oneTimeToken)
	if err == sql.ErrNoRows {
		return nil
	}
	var tokenTime time.Time
	er := DB.Get(&tokenTime, "SELECT onetime_login_token_at FROM users WHERE id=$1;", user.ID)
	if er == sql.ErrNoRows {
		return nil
	}
	if tokenTime.Add(time.Hour).Before(time.Now()) {
		return nil
	}
	_, e := DB.Exec("UPDATE users SET onetime_login_token_at = (datetime(0, 'unixepoch')) WHERE id=$1;", user.ID)
	if e != nil {
		glog.Errorf("Error resetting onetime sign-in token creation time: %s", e)
	}
	return &user
}

func UpdateUserPasswdByID(id uuid.UUID, passwd string) error {
	passwdHash := hashPassword(passwd)
	_, err := DB.Exec(`UPDATE users SET passwd_hash = $1 WHERE id=$2;`, passwdHash, id)
	if err != nil {
		glog.Errorf("Error updating password: %s", err.Error())
	}
	return err
}

func UpdateUserOneTimeLoginTokenByID(id uuid.UUID) string {
	token := uuid.New().String()
	_, err := DB.Exec(`UPDATE users SET onetime_login_token = $1, onetime_login_token_at = current_timestamp WHERE id=$2;`, token, id)
	if err != nil {
		glog.Errorf("Error updating one-time signin token: %s", err.Error())
	}
	return token
}

func LogOutUserByID(id uuid.UUID) error {
	_, err := DB.Exec(`UPDATE users SET logout_at = current_timestamp WHERE id=$1;`, id)
	return err
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		glog.Fatalf("Error hashing password: %s", err.Error())
	}
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
