// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"time"

	"github.com/golang/glog"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserID                      int          `db:"user_id"`
	DomainID                    int          `db:"domain_id"`
	Email                       string       `db:"email"`
	DisplayName                 string       `db:"display_name"`
	PasswdHash                  string       `db:"passwd_hash"`
	About                       string       `db:"about"`
	IsSuperAdmin                bool         `db:"is_superadmin"`
	IsSuperMod                  bool         `db:"is_supermod"`
	IsTopicAutoSub              bool         `db:"is_topic_autosubscribe"`
	IsCommentAutoSub            bool         `db:"is_comment_autosubscribe"`
	IsEmailNotificationDisabled bool         `db:"is_email_notifications_disabled"`
	NumTopics                   int          `db:"num_topics"`
	NumComments                 int          `db:"num_comments"`
	NumActivity                 int          `db:"num_activity"`
	OnetimeLoginToken           string       `db:"onetime_login_token"`
	OnetimeLoginTokenAt         sql.NullTime `db:"onetime_login_token_at"`
	ResetToken                  string       `db:"reset_token"`
	LastIP                      string       `db:"last_ip"`
	ActivityAt                  sql.NullTime `db:"activity_at"`
	ResetAt                     sql.NullTime `db:"reset_at"`
	LogoutAt                    sql.NullTime `db:"logout_at"`
	BannedAt                    sql.NullTime `db:"banned_at"`
	ArchivedAt                  sql.NullTime `db:"archived_at"`
	CreatedAt                   sql.NullTime `db:"created_at"`
	UpdatedAt                   sql.NullTime `db:"updated_at"`
}

func createUser(domainID int, email string, displayName string, passwd string, isSuperUser bool) error {
	passwdHash := hashPassword(passwd)
	_, err := DB.Exec(`INSERT INTO users(domain_id, email, display_name, passwd_hash, is_superadmin) VALUES($1, $2, $3, $4, $5);`,
		domainID, email, displayName, passwdHash, isSuperUser,
	)
	return err
}

func CreateUser(domainID int, email string, displayName string, passwd string) error {
	return createUser(domainID, email, displayName, passwd, false)
}

func CreateSuperUser(domainID int, email string, displayName string, passwd string) error {
	return createUser(domainID, email, displayName, passwd, true)
}

func ChangePasswd(domainID int, email string, passwd string) error {
	passwdHash := hashPassword(passwd)
	_, err := DB.Exec(`UPDATE users SET passwd_hash = $1 WHERE domain_id = $2 AND email = $3;`,
		passwdHash, domainID, email,
	)
	return err
}

func GetUserByID(userID int) *User {
	user := User{}
	err := DB.Get(&user, "SELECT * FROM users WHERE user_id=$1;", userID)
	if err != nil {
		if err != sql.ErrNoRows {
			glog.Errorf("Error reading user by email: %s\n", err.Error())
		}
		return nil
	}
	return &user
}

func GetUserByEmail(domainID int, email string) *User {
	user := User{}
	err := DB.Get(&user, "SELECT * FROM users WHERE domain_id=$1 AND email=$2;", domainID, email)
	if err != nil {
		if err != sql.ErrNoRows {
			glog.Errorf("Error reading user by email: %s\n", err.Error())
		}
		return nil
	}
	return &user
}

func GetUserByPasswd(domainID int, email string, passwd string) *User {
	user := User{}
	err := DB.Get(&user, "SELECT * FROM users WHERE domain_id=$1 AND email=$2;", domainID, email)
	if err != nil {
		if err != sql.ErrNoRows {
			glog.Errorf("Error getting user by passwd: %s", err.Error())
		}
		return nil
	}
	if !checkPasswordHash(passwd, user.PasswdHash) {
		return nil
	}
	return &user
}

func GetUserByOneTimeToken(domainID int, oneTimeToken string) *User {
	user := User{}
	err := DB.Get(&user, "SELECT * FROM users WHERE domain_id=$1 AND onetime_login_token=$1;", domainID, oneTimeToken)
	if err == sql.ErrNoRows {
		return nil
	}
	if !user.OnetimeLoginTokenAt.Valid || user.OnetimeLoginTokenAt.Time.Add(time.Hour).Before(time.Now()) {
		return nil
	}
	_, e := DB.Exec("UPDATE users SET onetime_login_token_at = to_timestamp(0) WHERE user_id=$1;", user.UserID)
	if e != nil {
		glog.Errorf("Error resetting onetime sign-in token creation time: %s", e)
	}
	return &user
}

func GetSuperModsByDomainID(domainID int) []User {
	users := []User{}
	DB.Select(&users, "SELECT * FROM users WHERE domain_id = $1 AND is_supermod = $2;", domainID, true)
	return users
}

func UpdateUserPasswdByID(userID int, passwd string) error {
	passwdHash := hashPassword(passwd)
	_, err := DB.Exec(`UPDATE users SET passwd_hash = $1 WHERE user_id=$2;`, passwdHash, userID)
	if err != nil {
		glog.Errorf("Error updating password: %s", err.Error())
	}
	return err
}

func UpdateUserOneTimeLoginTokenByID(userID int) string {
	token := uuid.New().String()
	_, err := DB.Exec(`UPDATE users SET onetime_login_token = $1, onetime_login_token_at = current_timestamp WHERE user_id=$2;`, token, userID)
	if err != nil {
		glog.Errorf("Error updating one-time signin token: %s", err.Error())
	}
	return token
}

func UpdateUserSuperUser(userID int, isSuperUser bool) {
	_, err := DB.Exec(`UPDATE users SET is_superadmin = $1 WHERE user_id=$2;`, isSuperUser, userID)
	if err != nil {
		glog.Errorf("Error updating superuser status: %s", err.Error())
	}
}

func UpdateUserSuperMod(userID int, isSuperMod bool) {
	_, err := DB.Exec(`UPDATE users SET is_supermod = $1 WHERE user_id=$2;`, isSuperMod, userID)
	if err != nil {
		glog.Errorf("Error updating supermod status: %s", err.Error())
	}
}

func LogOutUserByID(userID int) error {
	_, err := DB.Exec(`UPDATE users SET logout_at = current_timestamp WHERE user_id=$1;`, userID)
	return err
}

func DeleteUserByID(userID int) {
	_, err := DB.Exec(`DELETE FROM users WHERE user_id = $1;`, userID)
	if err != nil {
		glog.Errorf("Error deleting user: %s", err.Error())
	}
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
