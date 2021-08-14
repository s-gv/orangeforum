// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"errors"
	"math/rand"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

type Domain struct {
	DomainID               int          `db:"domain_id"`
	DomainName             string       `db:"domain_name"`
	ForumName              string       `db:"forum_name"`
	NoRegularSignupMsg     string       `db:"no_regular_signup_msg"`
	SignupToken            string       `db:"signup_token"`
	EditWindow             int          `db:"edit_window"`
	AutoTopicCloseDays     int          `db:"auto_topic_close_days"`
	UserActivityWindow     int          `db:"user_activity_window"`
	MaxNumActivity         int          `db:"max_num_activity"`
	IsRegularSignupEnabled bool         `db:"is_regular_signup_enabled"`
	IsReadOnly             bool         `db:"is_readonly"`
	IsGroupSub             bool         `db:"enable_group_sub"`
	IsTopicAutoSub         bool         `db:"enable_topic_autosub"`
	IsCommentAutoSub       bool         `db:"enable_comment_autosub"`
	ArchivedAt             sql.NullTime `db:"archived_at"`
	CreatedAt              sql.NullTime `db:"created_at"`
	UpdatedAt              sql.NullTime `db:"updated_at"`
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
	isRegularSignupEnabled bool,
	isReadOnly bool,
	signupToken string,
) {
	_, err := DB.Exec(`
	UPDATE domains
	SET
		forum_name = $2,
		is_regular_signup_enabled = $3,
		is_readonly = $4,
		signup_token = $5
	WHERE 
		domain_id = $1;`,
		domainID, forumName, isRegularSignupEnabled, isReadOnly, signupToken,
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
