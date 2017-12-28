// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/models/db"
	"github.com/s-gv/orangeforum/static"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup
	db.Init("sqlite3", ":memory:")
	models.Migrate()

	models.CreateSuperUser("admin", "admin12345")

	// Run tests
	retCode := m.Run()

	os.Exit(retCode)
}

func TestStyleHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/static/css/orangeforum.css", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(StyleHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	if body := rr.Body.String(); body != static.StyleSrc {
		t.Errorf("handler returned unexpected body: got %s", body)
	}

	if ccHeader := rr.Header().Get("Cache-Control"); ccHeader == "" {
		t.Errorf("Cache-Control header not set.")
	}

	if ctHeader := rr.Header().Get("Content-Type"); ctHeader != "text/css" {
		t.Errorf("Content-Type header incorrect. Got: %s\n", ctHeader)
	}
}

func TestIndexHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(IndexHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v", status)
	}

	if body := rr.Body.String(); !strings.Contains(body, "<a href=\"/login?next=%2F\">") {
		t.Errorf("Index page does not have link to the login page.")
	}
}
