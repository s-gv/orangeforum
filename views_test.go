// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/s-gv/orangeforum/models"
)

func getHTTPStr(url string) (err error, statusCode int, body string) {
	resp, err := http.DefaultClient.Get(TestServer.URL + url)
	if err != nil {
		return err, -1, ""
	}
	bodyBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return err2, -2, ""
	}
	return nil, resp.StatusCode, string(bodyBytes)
}

func getHTTPOKStr(url string) (err error, body string) {
	println(TestServer.URL)
	resp, err := http.DefaultClient.Get(TestServer.URL + url)
	if err != nil {
		return err, ""
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Got status code " + strconv.Itoa(resp.StatusCode) + ". was expecting " + strconv.Itoa(http.StatusOK)), ""
	}
	bodyBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return err2, ""
	}
	return nil, string(bodyBytes)
}

func TestDomainIndexPage(t *testing.T) {
	models.CleanDB()
	testDomainName := "test.com"
	if err := models.CreateDomain(testDomainName); err != nil {
		t.Errorf("Error creating domains: %s\n", err.Error())
	}

	err, body := getHTTPOKStr("/forums/" + testDomainName)
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
