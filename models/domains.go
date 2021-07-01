// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"errors"
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
	_, err := DB.Exec(`INSERT INTO domains(domain_name) VALUES($1);`, domainName)
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
