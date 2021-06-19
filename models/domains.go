// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

type Domain struct {
	DomainID   uuid.UUID `db:"domain_id"`
	DomainName string    `db:"domain_name"`
	ForumName  string    `db:"forum_name"`
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
