// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"regexp"
	"errors"
	"net/url"
)

func grabCSRFToken(body string) (string, error) {
	csrfToken := ""
	r := regexp.MustCompile("<input type=\"hidden\" name=\"csrf\" value=\"([A-Za-z0-9]+)\">")
	match := r.FindStringSubmatch(body)
	if len(match) > 0 {
		csrfToken = match[1]
	}
	if csrfToken == "" {
		return "", errors.New("Unable to find CSRF token")
	}
	return csrfToken, nil
}

func grabSessionID(recorder *httptest.ResponseRecorder) (string, error) {
	sessionid := ""
	r := regexp.MustCompile("^sessionid=([A-Za-z0-9]+);")
	for _, cookie := range recorder.HeaderMap["Set-Cookie"] {
		matches := r.FindStringSubmatch(cookie)
		if len(matches) > 0 {
			sessionid = matches[1]
		}
	}
	if sessionid == "" {
		return "", errors.New("Unable to find sessionid")
	}
	return sessionid, nil
}

func loginForTest(username string, passwd string) (string, error) {
	// GET login page
	loginReq, _ := http.NewRequest("GET", "/login", nil)
	loginRR := httptest.NewRecorder()
	http.HandlerFunc(LoginHandler).ServeHTTP(loginRR, loginReq)

	sessionid, serr := grabSessionID(loginRR)
	if serr != nil {
		return "", serr
	}

	csrfToken, cerr := grabCSRFToken(loginRR.Body.String())
	if cerr != nil {
		return "", cerr
	}

	// POST login and record sessionid
	form := url.Values{}
	form.Add("username", username)
	form.Add("passwd", passwd)
	form.Add("csrf", csrfToken)
	loginReq, _ = http.NewRequest("POST", "/login", strings.NewReader(form.Encode()))
	loginReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	loginReq.AddCookie(&http.Cookie{Name: "sessionid", Path: "/", Value: sessionid, HttpOnly: true})
	loginRR = httptest.NewRecorder()
	http.HandlerFunc(LoginHandler).ServeHTTP(loginRR, loginReq)

	if header, ok := loginRR.HeaderMap["Location"]; ok {
		if header[0] == "/" {
			return sessionid, nil
		} else {
			return "", errors.New("Unexpected re-direct after posting login. Maybe wrong password?")
		}
	}

	return "", errors.New("Login failed")
}

func TestUserProfileHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/users?u=admin", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserProfileHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "admin")  {
		t.Errorf("Profile page does not have the name of the user. Body: %s\n", body)
	}
	if strings.Contains(body, "private") {
		t.Errorf("Public profile page shows private info. Body: %s\n", body)
	}
}

func TestLoggedInUserProfileHandler(t *testing.T) {
	// Get the Login page and record CSRF token
	sessionid, err := loginForTest("admin", "admin12345")
	if err != nil {
		t.Fatalf("%v\n", err.Error())
	}

	req, _ := http.NewRequest("GET", "/users?u=admin", nil)
	req.AddCookie(&http.Cookie{Name: "sessionid", Path: "/", Value: sessionid, HttpOnly: true})
	rr := httptest.NewRecorder()
	http.HandlerFunc(UserProfileHandler).ServeHTTP(rr, req)

	body := rr.Body.String()
	if !strings.Contains(body, "admin")  {
		t.Errorf("Profile page does not have the name of the user. Body: %s\n", body)
	}
	if !strings.Contains(body, "private") {
		t.Errorf("Profile page doesn't show private info when logged in. Body: %s\n", body)
	}
	if !strings.Contains(body, "/changepass") {
		t.Errorf("Profile page doesn't have a link to change password page when logged in. Body: %s\n", body)
	}
}