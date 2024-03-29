// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"errors"
	"html/template"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
)

type Domain struct {
	DomainID                               int          `db:"domain_id"`
	DomainName                             string       `db:"domain_name"`
	ForumName                              string       `db:"forum_name"`
	NoRegularSignupMsg                     string       `db:"no_regular_signup_msg"`
	WhitelistedEmailDomains                string       `db:"whitelisted_email_domains"`
	SignupToken                            string       `db:"signup_token"`
	EditWindow                             int          `db:"edit_window"`
	AutoTopicCloseDays                     int          `db:"auto_topic_close_days"`
	UserActivityWindow                     int          `db:"user_activity_window"`
	MaxNumActivity                         int          `db:"max_num_activity"`
	HeaderMsg                              string       `db:"header_msg"`
	Logo                                   template.URL `db:"logo"`
	Icon                                   template.URL `db:"icon"`
	SMTPHost                               string       `db:"smtp_host"`
	SMTPPort                               int          `db:"smtp_port"`
	SMTPUser                               string       `db:"smtp_user"`
	SMTPPass                               string       `db:"smtp_pass"`
	DefaultFromEmail                       string       `db:"default_from_email"`
	IsPrivate                              bool         `db:"is_private"`
	IsRegularSigninEnabled                 bool         `db:"is_regular_signin_enabled"`
	IsAutoUserCreationOnEmailSigninEnabled bool         `db:"is_auto_user_creation_on_email_signin_enabled"`
	IsRegularSignupEnabled                 bool         `db:"is_regular_signup_enabled"`
	IsReadOnly                             bool         `db:"is_readonly"`
	IsGroupSub                             bool         `db:"enable_group_sub"`
	IsTopicAutoSub                         bool         `db:"enable_topic_autosub"`
	IsCommentAutoSub                       bool         `db:"enable_comment_autosub"`
	ArchivedAt                             sql.NullTime `db:"archived_at"`
	CreatedAt                              time.Time    `db:"created_at"`
	UpdatedAt                              time.Time    `db:"updated_at"`
}

func CreateDomain(domainName string) error {
	if len(domainName) < 2 {
		return errors.New("Domain name should have at least 2 characters")
	}
	if strings.Contains(domainName, " ") {
		return errors.New("Domain name should not have spaces")
	}
	signupToken := strconv.Itoa(rand.Intn(10000000) + 3245714)
	_, err := DB.Exec(`INSERT INTO domains(domain_name, signup_token) VALUES($1, $2);`, domainName, signupToken)
	return err
}

func GetDomainByName(domainName string) *Domain {
	var domain Domain
	err := DB.Get(&domain, `SELECT * FROM domains WHERE domain_name = $1;`, domainName)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		glog.Errorf("Error reading domain ID: %s", err.Error())
		return nil
	}
	return &domain
}

func GetDomainIDByName(domainName string) *int {
	var domainID int
	err := DB.Get(&domainID, `SELECT domain_id FROM domains WHERE domain_name = $1;`, domainName)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		glog.Errorf("Error reading domain ID: %s", err.Error())
		return nil
	}
	return &domainID
}

func UpdateDomainByID(
	domainID int,
	forumName string,
	whitelistedEmailDomains string,
	logo string,
	icon string,
	isRegularSignupEnabled bool,
	isRegularSigninEnabled bool,
	isAutoUserCreationOnEmailSigninEnabled bool,
	isReadOnly bool,
	isPrivate bool,
	signupToken string,
) {
	_, err := DB.Exec(`
	UPDATE domains
	SET
		forum_name = $2,
		whitelisted_email_domains = $3,
		is_regular_signup_enabled = $4,
		is_regular_signin_enabled = $5,
		is_auto_user_creation_on_email_signin_enabled = $6,
		is_readonly = $7,
		is_private = $8,
		signup_token = $9,
		logo = $10,
        icon = $11
	WHERE
		domain_id = $1;`,
		domainID,
		forumName,
		whitelistedEmailDomains,
		isRegularSignupEnabled,
		isRegularSigninEnabled,
		isAutoUserCreationOnEmailSigninEnabled,
		isReadOnly,
		isPrivate,
		signupToken,
		logo,
		icon,
	)
	if err != nil {
		glog.Errorf("Error updating domain ID:%d -- %s", domainID, err.Error())
	}
}

func UpdateDomainSMTPByID(domainID int, smtpHost string, smtpPort int, smtpUser string, smtpPass string, fromEmail string) {
	_, err := DB.Exec(`
	UPDATE domains
	SET
		smtp_host = $2,
		smtp_port = $3,
		smtp_user = $4,
		smtp_pass = $5,
		default_from_email = $6
	WHERE
		domain_id = $1;`,
		domainID, smtpHost, smtpPort, smtpUser, smtpPass, fromEmail,
	)
	if err != nil {
		glog.Errorf("Error updating domain ID:%d -- %s", domainID, err.Error())
	}
}

func DeleteDomainByID(domainID int) {
	_, err := DB.Exec(`DELETE FROM domains WHERE domain_id = $1;`, domainID)
	if err != nil {
		glog.Errorf("Error deleting domain ID:%d -- %s", domainID, err.Error())
	}
}
