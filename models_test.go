// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/s-gv/orangeforum/models"
)

func TestMain(m *testing.M) {
	models.DB = sqlx.MustConnect("pgx", "postgres://dbuser:dbpass@localhost:5432/oftestdb")
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestDomainCreation(t *testing.T) {
	models.CleanDB()
	testDomainName := "test.com"
	err := models.CreateDomain(testDomainName)
	if err != nil {
		t.Errorf("Error creating domains: %s\n", err.Error())
	}
	domain := models.GetDomainByName(testDomainName)
	if domain == nil {
		t.Errorf("Error getting domain\n")
	}
	if domain.DomainName != testDomainName {
		t.Errorf("Error getting domain. got: %s, expected: %s\n", domain.DomainName, testDomainName)
	}
}
