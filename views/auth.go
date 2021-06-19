// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"context"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/golang/glog"
	"github.com/gorilla/csrf"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/templates"
)

const CtxUserKey = contextKey("user")

var userNameReg *regexp.Regexp
var nextURLReg *regexp.Regexp
var emailReg *regexp.Regexp

func cleanNextURL(next string) string {
	if next == "" || next[0] != '/' {
		return "/"
	}
	return next
}

func authenticate(id int, w http.ResponseWriter) error {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"user_id": strconv.Itoa(id),
		"iat":     time.Now(),
		"exp":     time.Now().Add(365 * 24 * time.Hour),
	})
	if err == nil {
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Path:     "/",
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	}
	return err
}

func mustAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())
		basePath, _ := r.Context().Value(BasePath).(string)

		if err == nil && token != nil && jwt.Validate(token) == nil {
			if uid, ok := claims["user_id"].(string); ok {
				if iat, ok := claims["iat"].(time.Time); ok {
					userID, _ := strconv.Atoi(uid)
					user := models.GetUserByID(userID)
					if user != nil && user.LogoutAt.Time.Before(iat) {
						ctx := context.WithValue(r.Context(), CtxUserKey, user)
						// Token is authenticated, pass it through
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "jwt",
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-300 * time.Hour),
			HttpOnly: true,
		})
		if r.Method == "GET" {
			http.Redirect(w, r, basePath+"/auth/signin?next="+r.URL.Path, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, basePath+"/", http.StatusSeeOther)
		}
	})
}

func canAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err == nil && token != nil && jwt.Validate(token) == nil {
			if uid, ok := claims["user_id"].(string); ok {
				if iat, ok := claims["iat"].(time.Time); ok {
					userID, _ := strconv.Atoi(uid)
					user := models.GetUserByID(userID)
					if user != nil && user.LogoutAt.Time.Before(iat) {
						ctx := context.WithValue(r.Context(), CtxUserKey, user)
						// Token is authenticated, pass it through
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func getAuthSignIn(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(BasePath).(string)
	next := cleanNextURL(r.FormValue("next"))
	templates.Signin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"BasePath":       basePath,
		"Next":           next,
	})
}

func postAuthSignIn(w http.ResponseWriter, r *http.Request) {
	domainID := r.Context().Value(DomainID).(int)
	basePath := r.Context().Value(BasePath).(string)
	next := cleanNextURL(r.FormValue("next"))
	email := r.PostFormValue("email")
	passwd := r.PostFormValue("password")
	user := models.GetUserByPasswd(domainID, email, passwd)
	if user != nil {
		err := authenticate(user.UserID, w)
		if err != nil {
			glog.Errorf("Error authenticating: %s", err.Error())
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
	}
	templates.Signin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"BasePath":       basePath,
		"Next":           next,
		"ErrMsg":         "Invalid username / password",
	})
}

func getAuthOneTimeSignIn(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(BasePath).(string)
	next := cleanNextURL(r.FormValue("next"))
	templates.OneTimeSignin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"BasePath":       basePath,
		"Next":           next,
	})
}

func postAuthOneTimeSignIn(w http.ResponseWriter, r *http.Request) {
	domainID := r.Context().Value(DomainID).(int)
	basePath := r.Context().Value(BasePath).(string)
	next := cleanNextURL(r.PostFormValue("next"))
	email := r.PostFormValue("email")
	errMsg := "E-mail not found"

	user := models.GetUserByEmail(domainID, email)
	if user != nil {
		println(user.Email)
		errMsg = "A one time sign-in link has been sent to your email"
		token := models.UpdateUserOneTimeLoginTokenByID(user.UserID)
		link := "http://" + r.Host + basePath + "/auth/otsignin/" + token + "?next=" + next

		// domain := models.GetDomainByID(domainID) // TODO
		forumName := "Orange Forum" // TODO
		subject := forumName + " sign-in link"
		body := "Someone (hopefully you) requested a sign-in link for " + forumName + ".\r\n" +
			"If you want to sign-in, visit " + link + "\r\n\r\nIf not, just ignore this message."
		sendMail(user.Email, subject, body)
	}

	templates.OneTimeSignin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"BasePath":       basePath,
		"Next":           next,
		"ErrMsg":         errMsg,
	})
}

func getAuthOneTimeSignInDone(w http.ResponseWriter, r *http.Request) {
	domainID := r.Context().Value(DomainID).(int)
	next := cleanNextURL(r.FormValue("next"))
	token := chi.URLParam(r, "token")
	user := models.GetUserByOneTimeToken(domainID, token)
	if user != nil {
		if err := authenticate(user.UserID, w); err == nil {
			http.Redirect(w, r, next, http.StatusSeeOther)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func getAuthSignUp(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(BasePath).(string)
	next := cleanNextURL(r.FormValue("next"))
	templates.Signup.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"BasePath":       basePath,
		"Next":           next,
	})
}

func postAuthSignUp(w http.ResponseWriter, r *http.Request) {
	domainID := r.Context().Value(DomainID).(int)
	basePath := r.Context().Value(BasePath).(string)
	next := cleanNextURL(r.FormValue("next"))

	email := r.PostFormValue("email")
	passwd := r.PostFormValue("password")
	passwd2 := r.PostFormValue("password2")

	email = strings.Trim(email, " ")

	errMsg := ""
	if !strings.Contains(email, "@") {
		errMsg = "Invalid email"
	}
	if len(email) < 3 {
		errMsg = "Email should have at least 3 characters"
	} else if emailReg.ReplaceAllString(email, "") != email {
		errMsg = "Email should not have non-alphanumeric characters"
	} else if len(passwd) < 6 {
		errMsg = "Password should have at least 6 characters"
	} else if passwd != passwd2 {
		errMsg = "Passwords do not match"
	}

	existingUser := models.GetUserByEmail(domainID, email)
	if existingUser != nil {
		errMsg = "E-mail already registered"
	}

	if errMsg == "" {
		userName := strings.Split(email, "@")[0] + strconv.Itoa(rand.Intn(100000000))

		err := models.CreateUser(domainID, email, userName, passwd)
		if err != nil {
			glog.Errorf("Error creating user: %s", err.Error())
			errMsg = "Error during signup."
		}
	}

	if errMsg == "" {
		glog.Infof("Created user: %s for domainID: %d", email, domainID)
		user := models.GetUserByEmail(domainID, email)
		err := authenticate(user.UserID, w)
		if err != nil {
			glog.Errorf("Error authenticating: %s", err.Error())
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
	}

	templates.Signup.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"BasePath":       basePath,
		"ErrMsg":         errMsg,
		"Next":           next,
	})
}

func getAuthChangePass(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(BasePath).(string)
	templates.ChangePass.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"BasePath":       basePath,
	})
}

func postAuthChangePass(w http.ResponseWriter, r *http.Request) {
	domainID := r.Context().Value(DomainID).(int)
	basePath := r.Context().Value(BasePath).(string)

	user := r.Context().Value(CtxUserKey).(*models.User)
	oldPasswd := r.PostFormValue("old_password")
	passwd := r.PostFormValue("password")
	passwd2 := r.PostFormValue("password2")

	var errMsg string
	if len(passwd) < 6 {
		errMsg = "New password should have at least 6 characters"
	} else if passwd != passwd2 {
		errMsg = "New passwords do not match"
	}

	oldUser := models.GetUserByPasswd(domainID, user.Email, oldPasswd)
	if oldUser == nil || oldUser.UserID != user.UserID {
		errMsg = "Old password incorrect"
	}

	if errMsg == "" {
		err := models.UpdateUserPasswdByID(user.UserID, passwd)
		if err == nil {
			errMsg = "Changed password successfully"
		} else {
			errMsg = "Error changing password"
			glog.Errorf("Error changing password: %s", err.Error())
		}
	}
	templates.ChangePass.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"BasePath":       basePath,
		"ErrMsg":         errMsg,
	})
}

func getAuthLogout(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(BasePath).(string)

	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: "", Path: "/", Expires: time.Now().Add(-300 * time.Hour), HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "csrftoken", Value: "", Path: "/", Expires: time.Now().Add(-300 * time.Hour)})
	if user, ok := r.Context().Value(CtxUserKey).(*models.User); ok {
		models.LogOutUserByID(user.UserID)
	}
	http.Redirect(w, r, basePath+"/", http.StatusSeeOther)
}

func init() {
	userNameReg = regexp.MustCompile("[^a-zA-Z0-9]+")
	emailReg = regexp.MustCompile("[^a-zA-Z0-9@\\.-]+")
	nextURLReg = regexp.MustCompile("[^a-zA-Z0-9-/]+")
}
