// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/s-gv/orangeforum/models"
)

const (
	testDomainName = "test.com"
	testAdminEmail = "admin@example.com"
	testAdminName  = "Admin User"
	testAdminPass  = "testpass123"
)

func getHTTPOKStr(c *http.Client, url string) (err error, body string) {
	resp, err := c.Get(TestServer.URL + url)
	if err != nil {
		return err, ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Got status code " + strconv.Itoa(resp.StatusCode) + ". was expecting " + strconv.Itoa(http.StatusOK)), ""
	}
	bodyBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return err2, ""
	}
	return nil, string(bodyBytes)
}

func loginAs(c *http.Client, domainName string, email string, password string) error {
	resp, err := c.PostForm(TestServer.URL+"/forums/"+domainName+"/auth/signin", url.Values{
		"email": {email}, "password": {password}})
	if err != nil {
		return errors.New("Error posting to /auth/signin")
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return errors.New("Error reading body after posting to /auth/signin")
	}

	body := string(bodyBytes)
	if !strings.Contains(body, "Logout") {
		return errors.New("Signing in failed.")
	}
	return nil
}

func createTestDomainAndUsers() error {
	if err := models.CreateDomain(testDomainName); err != nil {
		return err
	}
	domain := models.GetDomainByName(testDomainName)
	if domain == nil {
		return errors.New("Error reading domain: " + testDomainName)
	}
	if err := models.CreateSuperUser(domain.DomainID, testAdminEmail, testAdminName, testAdminPass); err != nil {
		return errors.New("Error creating admin: " + testAdminEmail)
	}
	return nil
}

func TestDomainIndexPage(t *testing.T) {
	models.CleanDB()

	if err := models.CreateDomain(testDomainName); err != nil {
		t.Errorf("Error creating domains: %s\n", err.Error())
	}

	err, body := getHTTPOKStr(&http.Client{}, "/forums/"+testDomainName)
	if err != nil {
		t.Errorf("Error getting index page: %s\n", err.Error())
	}
	if !strings.Contains(body, "Login") {
		t.Errorf("Expected to see the Login button on the index page\n")
	}
	if !strings.Contains(body, "Signup") {
		t.Errorf("Expected to see the Signup button on the index page\n")
	}
}

func TestAuthedDomainIndexPage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	if err := loginAs(client, testDomainName, testAdminEmail, testAdminPass); err != nil {
		t.Errorf("Error signing in: %s\n", err.Error())
	}

	err, body := getHTTPOKStr(client, "/forums/"+testDomainName+"/")
	if err != nil {
		t.Errorf("Error getting index page: %s\n", err.Error())
	}
	if !strings.Contains(body, testAdminName) {
		t.Errorf("Index page does not contain the display name.\n")
	}
}
