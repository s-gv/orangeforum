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
	testUser2Email = "user2@example.com"
	testUser2Name  = "Jane Doe"
	testUser2Pass  = "testuser2pass123"
)

func getHTTPOKStr(c *http.Client, url string) (body string, err error) {
	resp, err := c.Get(TestServer.URL + url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Got status code " + strconv.Itoa(resp.StatusCode) + ". was expecting " + strconv.Itoa(http.StatusOK))
	}
	bodyBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return "", err2
	}
	return string(bodyBytes), nil
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

func postHTTPOKStr(c *http.Client, url string, values url.Values) (body string, err error) {
	resp, err := c.PostForm(TestServer.URL+url, values)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Got status code " + strconv.Itoa(resp.StatusCode) + ". was expecting " + strconv.Itoa(http.StatusOK))
	}
	bodyBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return "", err2
	}
	return string(bodyBytes), nil
}

func postHTTPForbidden(c *http.Client, url string, values url.Values) error {
	resp, err := c.PostForm(TestServer.URL+url, values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusForbidden {
		return errors.New("Got status code " + strconv.Itoa(resp.StatusCode) + ". was expecting " + strconv.Itoa(http.StatusForbidden))
	}
	return nil
}

func loginAs(c *http.Client, domainName string, email string, password string) error {
	body, err := postHTTPOKStr(c, "/forums/"+domainName+"/auth/signin", url.Values{
		"email": {email}, "password": {password}})
	if err != nil {
		return errors.New("Error posting to /auth/signin")
	}
	if !strings.Contains(body, "Logout") {
		return errors.New("signing in failed")
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
	if err := models.CreateUser(domain.DomainID, testUser2Email, testUser2Name, testUser2Pass); err != nil {
		return errors.New("Error creating user: " + testUser2Email)
	}
	return nil
}

func TestDomainIndexPage(t *testing.T) {
	models.CleanDB()

	if err := models.CreateDomain(testDomainName); err != nil {
		t.Errorf("Error creating domains: %s\n", err.Error())
	}

	body, err := getHTTPOKStr(&http.Client{}, "/forums/"+testDomainName)
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

	body, err := getHTTPOKStr(client, "/forums/"+testDomainName+"/")
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

	body, err := getHTTPOKStr(client, "/forums/"+testDomainName+"/admin")
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

	body, err := postHTTPOKStr(client, "/forums/"+testDomainName+"/admin", url.Values{
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

func TestProfileUpdatePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}
	domain := models.GetDomainByName(testDomainName)
	user := models.GetUserByEmail(domain.DomainID, testUserEmail)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	body, err := getHTTPOKStr(client, "/forums/"+testDomainName+"/users/"+strconv.Itoa(user.UserID))
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(body, user.DisplayName) {
		t.Errorf("Expected to find display name in profile page: %s\n", user.DisplayName)
	}

	if err := postHTTPForbidden(
		client,
		"/forums/"+testDomainName+"/users/"+strconv.Itoa(user.UserID),
		url.Values{"display_name": {"Baby Doe"}, "email": {"doe@example.com"}},
	); err != nil {
		t.Error(err)
	}
}

func TestBadAuthProfileUpdatePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}
	domain := models.GetDomainByName(testDomainName)
	user := models.GetUserByEmail(domain.DomainID, testUserEmail)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testUser2Email, testUser2Pass)

	body, err := getHTTPOKStr(client, "/forums/"+testDomainName+"/users/"+strconv.Itoa(user.UserID))
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(body, user.DisplayName) {
		t.Errorf("Expected to find display name in profile page: %s\n", user.DisplayName)
	}

	if err := postHTTPForbidden(
		client,
		"/forums/"+testDomainName+"/users/"+strconv.Itoa(user.UserID),
		url.Values{"display_name": {"Baby Doe"}, "email": {"doe@example.com"}},
	); err != nil {
		t.Error(err)
	}
}

func TestAuthedProfileUpdatePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}
	domain := models.GetDomainByName(testDomainName)
	user := models.GetUserByEmail(domain.DomainID, testUserEmail)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testUserEmail, testUserPass)

	body, err := getHTTPOKStr(client, "/forums/"+testDomainName+"/users/"+strconv.Itoa(user.UserID))
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(body, user.DisplayName) {
		t.Errorf("Expected to find display name in profile page: %s\n", user.DisplayName)
	}

	newDisplayName := "Baby Doe"
	newEmail := "babydoe@example.com"
	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/users/"+strconv.Itoa(user.UserID),
		url.Values{"display_name": {newDisplayName}, "email": {newEmail}},
	); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(body, newDisplayName) {
			t.Errorf("Expected to find new display name in profile page: %s\n", newDisplayName)
		}
		if !strings.Contains(body, newEmail) {
			t.Errorf("Expected to find new email in profile page: %s\n", newEmail)
		}
	}
}

func TestAdminModCreatePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}
	domain := models.GetDomainByName(testDomainName)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testAdminEmail, testAdminPass)

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/admin/mods/create",
		url.Values{"mod_user_email": {testUserEmail}, "action": {"Add"}},
	); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(body, testUserEmail) {
			t.Errorf("Expected to find email of new mod in admin page: %s\n", testUserEmail)
		}
		if !strings.Contains(body, testUserName) {
			t.Errorf("Expected to find display name of new mod in admin page : %s\n", testUserName)
		}
		user := models.GetUserByEmail(domain.DomainID, testUserEmail)
		if !user.IsSuperMod {
			t.Errorf("Expected %s to be a mod.\n", user.Email)
		}
	}
}

func TestAdminModDeletePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}
	domain := models.GetDomainByName(testDomainName)
	user := models.GetUserByEmail(domain.DomainID, testUserEmail)

	models.UpdateUserSuperMod(user.UserID, true)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testAdminEmail, testAdminPass)

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/admin/mods/delete",
		url.Values{"mod_user_id": {strconv.Itoa(user.UserID)}, "action": {"Remove"}},
	); err != nil {
		t.Error(err)
	} else {
		if strings.Contains(body, testUserEmail) {
			t.Errorf("Expected to not find email of mod in admin page: %s\n", testUserEmail)
		}
		if strings.Contains(body, testUserName) {
			t.Errorf("Expected to not find display name of new mod in admin page : %s\n", testUserName)
		}

		if models.GetUserByEmail(domain.DomainID, testUserEmail).IsSuperMod {
			t.Errorf("Expected %s to not be a mod.\n", user.Email)
		}
	}
}

func TestAdminCategoryCreatePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testAdminEmail, testAdminPass)

	newCategoryName := "Off-topic"
	newCategoryDesc := "Topics that dont belong elsewhere"

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/admin/categories/create",
		url.Values{"name": {newCategoryName}, "description": {newCategoryDesc}, "action": {"Add"}},
	); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(body, newCategoryName) {
			t.Errorf("Expected to find name of new category in admin page: %s\n", newCategoryName)
		}
		if !strings.Contains(body, newCategoryDesc) {
			t.Errorf("Expected to find the description of the new category in admin page : %s\n", newCategoryDesc)
		}
	}

	if body, err := getHTTPOKStr(client, "/forums/"+testDomainName+"/"); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(body, newCategoryName) {
			t.Errorf("Expected to find name of new category in main page: %s\n", newCategoryName)
		}
		if !strings.Contains(body, newCategoryDesc) {
			t.Errorf("Expected to find the description of the new category in main page : %s\n", newCategoryDesc)
		}
	}
}

func TestAdminCategoryUpdatePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	domain := models.GetDomainByName(testDomainName)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testAdminEmail, testAdminPass)

	newCategoryName := "Off-topic"
	newCategoryDesc := "Topics that dont belong elsewhere"

	categoryID := models.CreateCategory(domain.DomainID, newCategoryName, newCategoryDesc)

	updatedCategoryName := "Politics"
	updatedCategoryDesc := "Topics that create flame wars"

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/admin/categories/"+strconv.Itoa(categoryID),
		url.Values{
			"name":        {updatedCategoryName},
			"description": {updatedCategoryDesc},
			"is_readonly": {"1"},
			"is_private":  {"1"},
			"is_archived": {"1"},
			"action":      {"Update"},
		},
	); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(body, updatedCategoryName) {
			t.Errorf("Expected to find name of new category in admin page: %s\n", updatedCategoryName)
		}
		if !strings.Contains(body, updatedCategoryDesc) {
			t.Errorf("Expected to find the description of the new category in admin page: %s\n", updatedCategoryDesc)
		}
	}

	updatedCategory := models.GetCategoryByID(categoryID)
	if !(updatedCategory.IsReadOnly && updatedCategory.IsPrivate && updatedCategory.ArchivedAt.Valid) {
		t.Errorf("Category not updated properly: %v\n", updatedCategory)
	}

	if body, err := getHTTPOKStr(client, "/forums/"+testDomainName+"/"); err != nil {
		t.Error(err)
	} else {
		if strings.Contains(body, updatedCategoryName) {
			t.Errorf("Expected to not find name of archived category in main page: %s\n", updatedCategoryName)
		}
		if strings.Contains(body, updatedCategoryDesc) {
			t.Errorf("Expected to not find the description of the archived category in main page: %s\n", updatedCategoryDesc)
		}
	}
}

func TestTopicCreationPage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	domain := models.GetDomainByName(testDomainName)
	newCategoryName := "Off-topic"
	newCategoryDesc := "Topics that dont belong elsewhere"
	newTopicName := "Osama killed in drone strike"
	newTopicContent := "Osama was found hiding in Pakistan"

	categoryID := models.CreateCategory(domain.DomainID, newCategoryName, newCategoryDesc)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testUserEmail, testUserPass)

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/categories/"+strconv.Itoa(categoryID)+"/topics/new",
		url.Values{
			"title":   {newTopicName},
			"content": {newTopicContent},
			"action":  {"Submit"},
		},
	); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(body, newTopicName) {
			t.Errorf("Expected to find name of new topic: %s\n", newTopicName)
		}
		if !strings.Contains(body, newTopicContent) {
			t.Errorf("Expected to find the content of new topic: %s\n", newTopicContent)
		}
	}
}

func TestTopicUpdatePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	domain := models.GetDomainByName(testDomainName)
	user := models.GetUserByEmail(domain.DomainID, testUserEmail)
	newCategoryName := "Off-topic"
	newCategoryDesc := "Topics that dont belong elsewhere"
	newTopicTitle := "Osama killed in drone strike"
	newTopicContent := "Osama was found hiding in Pakistan"

	categoryID := models.CreateCategory(domain.DomainID, newCategoryName, newCategoryDesc)
	topicID := models.CreateTopic(categoryID, user.UserID, newTopicTitle, newTopicContent, false)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testUserEmail, testUserPass)

	updatedTopicTitle := "Puppy is doing zoomies"
	updatedTopicContent := "The puppy is cute"

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/categories/"+strconv.Itoa(categoryID)+"/topics/"+strconv.Itoa(topicID)+"/edit",
		url.Values{
			"title":   {updatedTopicTitle},
			"content": {updatedTopicContent},
			"action":  {"Update"},
		},
	); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(body, updatedTopicTitle) {
			t.Errorf("Expected to find title of the topic: %s\n", updatedTopicTitle)
		}
		if !strings.Contains(body, updatedTopicContent) {
			t.Errorf("Expected to find the content of the topic: %s\n", updatedTopicContent)
		}
	}
}

func TestTopicDeletePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	domain := models.GetDomainByName(testDomainName)
	user := models.GetUserByEmail(domain.DomainID, testUserEmail)
	newCategoryName := "Off-topic"
	newCategoryDesc := "Topics that dont belong elsewhere"
	newTopicTitle := "Osama killed in drone strike"
	newTopicContent := "Osama was found hiding in Pakistan"

	categoryID := models.CreateCategory(domain.DomainID, newCategoryName, newCategoryDesc)
	topicID := models.CreateTopic(categoryID, user.UserID, newTopicTitle, newTopicContent, false)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testUserEmail, testUserPass)

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/categories/"+strconv.Itoa(categoryID)+"/topics/"+strconv.Itoa(topicID)+"/edit",
		url.Values{
			"title":   {newTopicTitle},
			"content": {newTopicContent},
			"action":  {"Delete"},
		},
	); err != nil {
		t.Error(err)
	} else {
		if strings.Contains(body, newTopicTitle) {
			t.Errorf("Expected to not find name of new topic: %s\n", newTopicTitle)
		}
		if strings.Contains(body, newTopicContent) {
			t.Errorf("Expected to not find the content of new topic: %s\n", newTopicContent)
		}
	}
}

func TestCommentCreatePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	domain := models.GetDomainByName(testDomainName)
	user := models.GetUserByEmail(domain.DomainID, testUserEmail)
	newCategoryName := "Off-topic"
	newCategoryDesc := "Topics that dont belong elsewhere"
	newTopicTitle := "Osama killed in drone strike"
	newTopicContent := "Osama was found hiding in Pakistan"
	newCommentContent := "Serves him right"

	categoryID := models.CreateCategory(domain.DomainID, newCategoryName, newCategoryDesc)
	topicID := models.CreateTopic(categoryID, user.UserID, newTopicTitle, newTopicContent, false)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testUserEmail, testUserPass)

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/categories/"+strconv.Itoa(categoryID)+"/topics/"+strconv.Itoa(topicID)+"/comments/new",
		url.Values{
			"content": {newCommentContent},
			"action":  {"Submit"},
		},
	); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(body, newCommentContent) {
			t.Errorf("Expected to find the content of new comment: %s\n", newCommentContent)
		}
	}
}

func TestCommentUpdatePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	domain := models.GetDomainByName(testDomainName)
	user := models.GetUserByEmail(domain.DomainID, testUserEmail)
	newCategoryName := "Off-topic"
	newCategoryDesc := "Topics that dont belong elsewhere"
	newTopicTitle := "Osama killed in drone strike"
	newTopicContent := "Osama was found hiding in Pakistan"
	newCommentContent := "Serves him right"

	categoryID := models.CreateCategory(domain.DomainID, newCategoryName, newCategoryDesc)
	topicID := models.CreateTopic(categoryID, user.UserID, newTopicTitle, newTopicContent, false)
	commentID := models.CreateComment(topicID, user.UserID, newCommentContent, false)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testUserEmail, testUserPass)

	updatedCommentContent := "How wonderful"

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/categories/"+strconv.Itoa(categoryID)+"/topics/"+strconv.Itoa(topicID)+"/comments/"+strconv.Itoa(commentID)+"/edit",
		url.Values{
			"content": {updatedCommentContent},
			"action":  {"Update"},
		},
	); err != nil {
		t.Error(err)
	} else {
		if !strings.Contains(body, updatedCommentContent) {
			t.Errorf("Expected to find the content of updated comment: %s\n", updatedCommentContent)
		}
	}
}

func TestCommentDeletePage(t *testing.T) {
	models.CleanDB()

	if err := createTestDomainAndUsers(); err != nil {
		t.Error(err)
	}

	domain := models.GetDomainByName(testDomainName)
	user := models.GetUserByEmail(domain.DomainID, testUserEmail)
	newCategoryName := "Off-topic"
	newCategoryDesc := "Topics that dont belong elsewhere"
	newTopicTitle := "Osama killed in drone strike"
	newTopicContent := "Osama was found hiding in Pakistan"
	newCommentContent := "Serves him right"

	categoryID := models.CreateCategory(domain.DomainID, newCategoryName, newCategoryDesc)
	topicID := models.CreateTopic(categoryID, user.UserID, newTopicTitle, newTopicContent, false)
	commentID := models.CreateComment(topicID, user.UserID, newCommentContent, false)

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	loginAs(client, testDomainName, testUserEmail, testUserPass)

	updatedCommentContent := "How wonderful"

	if body, err := postHTTPOKStr(
		client,
		"/forums/"+testDomainName+"/categories/"+strconv.Itoa(categoryID)+"/topics/"+strconv.Itoa(topicID)+"/comments/"+strconv.Itoa(commentID)+"/edit",
		url.Values{
			"content": {updatedCommentContent},
			"action":  {"Delete"},
		},
	); err != nil {
		t.Error(err)
	} else {
		if strings.Contains(body, updatedCommentContent) || strings.Contains(body, newCommentContent) {
			t.Errorf("Expected to not find the content of comment\n")
		}
	}
}
