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
	testUserEmail  = "user@example.com"
	testUserName   = "John Doe"
	testUserPass   = "testuserpass123"
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

func getHTTPForbidden(c *http.Client, url string) error {
	resp, err := c.Get(TestServer.URL + url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		return errors.New("Got status code " + strconv.Itoa(resp.StatusCode) + ". was expecting " + strconv.Itoa(http.StatusForbidden))
	}
	return nil
}

func postHTTPOKStr(c *http.Client, url string, values url.Values) (err error, body string) {
	resp, err := c.PostForm(TestServer.URL+url, values)
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
	err, body := postHTTPOKStr(c, "/forums/"+domainName+"/auth/signin", url.Values{
		"email": {email}, "password": {password}})
	if err != nil {
		return errors.New("Error posting to /auth/signin")
	}
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
	if err := models.CreateUser(domain.DomainID, testUserEmail, testUserName, testUserPass); err != nil {
		return errors.New("Error creating user: " + testUserEmail)
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

func TestAuthedAdminPage(t *testing.T) {
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

	err, body := getHTTPOKStr(client, "/forums/"+testDomainName+"/admin")
	if err != nil {
		t.Errorf("Error getting admin page: %s\n", err.Error())
	}
	if !strings.Contains(body, testAdminName) {
		t.Errorf("Admin page does not contain the display name.\n")
	}
}

func TestAdminWithoutPrivilegePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	if err := loginAs(client, testDomainName, testUserEmail, testUserPass); err != nil {
		t.Errorf("Error signing in: %s\n", err.Error())
	}

	err := getHTTPForbidden(client, "/forums/"+testDomainName+"/admin")
	if err != nil {
		t.Errorf("Admin page should not be accessible: %s\n", err.Error())
	}
}

func TestAuthedAdminUpdatePage(t *testing.T) {
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

	newForumName := "Orwell News"
	newIsRegularSignupEnabled := "1"
	newIsReadOnly := "1"

	err, body := postHTTPOKStr(client, "/forums/"+testDomainName+"/admin", url.Values{
		"forum_name":                {newForumName},
		"is_regular_signup_enabled": {newIsRegularSignupEnabled},
		"is_readonly":               {newIsReadOnly},
	})
	if err != nil {
		t.Errorf("Error updating admin page: %s\n", err.Error())
	}
	if !strings.Contains(body, newForumName) {
		t.Errorf("Expected new forum name %s in the returned page\n", newForumName)
	}

	domain := models.GetDomainByName(testDomainName)
	if domain == nil {
		t.Errorf("Error reading domain\n")
	}
	if domain != nil {
		if domain.ForumName != newForumName {
			t.Errorf("Expected forum name: %s, got: %s\n", newForumName, domain.ForumName)
		}
		if domain.IsRegularSignupEnabled != (newIsRegularSignupEnabled == "1") {
			t.Errorf("Expected IsRegularSignupEnabled: %s, got: %v\n", newIsRegularSignupEnabled, domain.IsRegularSignupEnabled)
		}
		if domain.IsReadOnly != (newIsReadOnly == "1") {
			t.Errorf("Expected IsReadOnly: %s, got: %v\n", newIsReadOnly, domain.IsReadOnly)
		}
	}
}
