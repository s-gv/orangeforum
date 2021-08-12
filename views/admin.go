// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"math/rand"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/templates"
)

func getAdmin(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	if !user.IsSuperAdmin {
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return
	}

	templates.Admin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		"Domain":         domain,
		"Host":           r.Host,
	})
}

func postAdmin(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	if !user.IsSuperAdmin {
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}

	errMsg := ""
	forumName := r.PostFormValue("forum_name")
	isRegularSignupEnabled := r.PostFormValue("is_regular_signup_enabled") == "1"
	isReadOnly := r.PostFormValue("is_readonly") == "1"
	signupToken := r.PostFormValue("signup_token")

	if !isRegularSignupEnabled && signupToken == "" {
		var letterRunes = []rune("0123456789")
		b := make([]rune, 12)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		signupToken = string(b)
	}

	if len(forumName) < 3 || len(forumName) > 30 {
		errMsg = "Forum name should have between 3 and 30 characters"
	}

	if len(signupToken) > 30 {
		errMsg = "Signup token should have fewer than 30 characters"
	}

	for _, r := range signupToken {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' {
			errMsg = "Signup token should be alphanumeric"
		}
	}

	if errMsg == "" {
		models.UpdateDomainByID(domain.DomainID, forumName, isRegularSignupEnabled, isReadOnly, signupToken)
		http.Redirect(w, r, basePath+"/admin", http.StatusSeeOther)
		return
	}

	templates.Admin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		"Domain":         domain,
		"Host":           r.Host,
		"ErrMsg":         errMsg,
	})
}