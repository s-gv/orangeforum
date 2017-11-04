// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

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

	if body := rr.Body.String(); !strings.Contains(body, "admin")  {
		t.Errorf("Profile page does not have the name of the user. Body: %s\n", body)
	}
}
