// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/templates"
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	if !user.IsSuperAdmin {
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		return
	}

	errMsg := ""

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form.", http.StatusBadRequest)
			return
		}

		forumName := r.PostFormValue("forum_name")
		logo := r.PostFormValue("logo")
		isRegularSignupEnabled := r.PostFormValue("is_regular_signup_enabled") == "1"
		isReadOnly := r.PostFormValue("is_readonly") == "1"
		signupToken := r.PostFormValue("signup_token")

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
			models.UpdateDomainByID(domain.DomainID, forumName, logo, isRegularSignupEnabled, isReadOnly, signupToken)
			http.Redirect(w, r, basePath+"/admin", http.StatusSeeOther)
			return
		}
	}

	mods := models.GetSuperModsByDomainID(domain.DomainID)
	categories := models.GetCategoriesByDomainID(domain.DomainID)

	templates.Admin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		UserField:        user,
		"Domain":         domain,
		"Host":           r.Host,
		"Mods":           mods,
		"Categories":     categories,
		"ErrMsg":         errMsg,
	})
}

func postCreateMod(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	if !user.IsSuperAdmin {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}

	modUserEmail := r.PostFormValue("mod_user_email")
	modUser := models.GetUserByEmail(domain.DomainID, modUserEmail)

	if modUser == nil {
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}

	models.UpdateUserSuperMod(modUser.UserID, true)

	http.Redirect(w, r, basePath+"/admin", http.StatusSeeOther)
}

func postDeleteMod(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	if !user.IsSuperAdmin {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}

	modUserID, err := strconv.Atoi(r.PostFormValue("mod_user_id"))
	if err != nil {
		http.Error(w, "Error reading user ID.", http.StatusBadRequest)
		return
	}
	modUser := models.GetUserByID(modUserID)

	if modUser.DomainID != domain.DomainID {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	models.UpdateUserSuperMod(modUserID, false)

	http.Redirect(w, r, basePath+"/admin", http.StatusSeeOther)
}

func postCategory(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	if domain.DomainID != user.DomainID || !user.IsSuperAdmin {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	isPrivate := r.PostFormValue("is_private") == "1"
	isReadOnly := r.PostFormValue("is_readonly") == "1"
	isArchived := r.PostFormValue("is_archived") == "1"
	categoryIDStr := chi.URLParam(r, "categoryID")

	if categoryIDStr == "create" {
		models.CreateCategory(domain.DomainID, name, description)
	} else {
		categoryID, err := strconv.Atoi(categoryIDStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		models.UpdateCategoryByID(categoryID, name, description, isPrivate, isReadOnly, isArchived)
	}

	http.Redirect(w, r, basePath+"/admin", http.StatusSeeOther)
}
