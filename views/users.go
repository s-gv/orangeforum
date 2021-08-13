// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/templates"
)

const UserField = "User"

func getProfile(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	profileUserID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Bad user ID", http.StatusBadRequest)
		return
	}
	profileUser := models.GetUserByID(profileUserID)
	if profileUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	showForm := (user.UserID == profileUser.UserID)
	if domain.DomainID == user.DomainID {
		if user.IsSuperAdmin || user.IsSuperMod {
			showForm = true
		}
	}

	showBan := (domain.DomainID == user.DomainID) && (user.IsSuperMod || user.IsSuperAdmin)

	templates.Profile.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		UserField:        user,
		"ProfileUser":    profileUser,
		"ShowForm":       showForm,
		"ShowBan":        showBan,
	})
}

func postProfile(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	user := r.Context().Value(CtxUserKey).(*models.User)
	domain, _ := r.Context().Value(ctxDomain).(*models.Domain)

	profileUserID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Bad user ID", http.StatusBadRequest)
		return
	}
	profileUser := models.GetUserByID(profileUserID)
	if profileUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	showForm := (user.UserID == profileUser.UserID)
	if domain.DomainID == user.DomainID {
		if user.IsSuperAdmin || user.IsSuperMod {
			showForm = true
		}
	}

	showBan := (domain.DomainID == user.DomainID) && (user.IsSuperMod || user.IsSuperAdmin)

	if profileUser.UserID != user.UserID {
		if (user.DomainID != domain.DomainID) || !(user.IsSuperMod || user.IsSuperAdmin) {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	newEmail := strings.Trim(r.PostFormValue("email"), " ")
	newDisplayName := r.PostFormValue("display_name")
	newIsBanned := r.PostFormValue("is_banned") == "1"

	if profileUser.IsSuperAdmin || profileUser.IsSuperMod {
		newIsBanned = false
	}

	errMsg := ""

	errMsg = validateEmail(newEmail)

	if len(newDisplayName) < 3 || len(newDisplayName) > 30 {
		errMsg = "Display name should have between 3 and 30 characters"
	}

	if errMsg == "" {
		models.UpdateUserByID(profileUser.UserID, newEmail, newDisplayName, newIsBanned)
		http.Redirect(w, r, basePath+"/users/"+strconv.Itoa(profileUser.UserID), http.StatusSeeOther)
		return
	}

	templates.Profile.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		UserField:        user,
		"ProfileUser":    profileUser,
		"ShowForm":       showForm,
		"ShowBan":        showBan,
		"ErrMsg":         errMsg,
	})
}
