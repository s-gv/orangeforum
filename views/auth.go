// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/templates"
)

const CtxUserKey string = "user"

var userNameReg *regexp.Regexp
var nextURLReg *regexp.Regexp

func cleanNextURL(next string) string {
	if next == "" || strings.Contains(next, ":") || next[0] != '/' || nextURLReg.ReplaceAllString(next, "") != next {
		return "/"
	}
	return next
}

func authenticate(id uuid.UUID, w http.ResponseWriter) error {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"user_id": id.String(),
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

		if err == nil && token != nil && jwt.Validate(token) == nil {
			if uid, ok := claims["user_id"].(string); ok {
				userID, er := uuid.Parse(uid)
				if er == nil {
					if iat, ok := claims["iat"].(time.Time); ok {
						user := models.GetUserByID(userID)
						if user != nil && user.LoggedOutAt.Before(iat) {
							ctx := context.WithValue(r.Context(), CtxUserKey, user)
							// Token is authenticated, pass it through
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						}
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
			http.Redirect(w, r, "/auth/signin?next="+r.URL.Path, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})
}

func canAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := jwtauth.FromContext(r.Context())

		if err == nil && token != nil && jwt.Validate(token) == nil {
			if uid, ok := claims["user_id"].(string); ok {
				userID, er := uuid.Parse(uid)
				if er == nil {
					if iat, ok := claims["iat"].(time.Time); ok {
						user := models.GetUserByID(userID)
						if user != nil && user.LoggedOutAt.Before(iat) {
							ctx := context.WithValue(r.Context(), CtxUserKey, user)
							// Token is authenticated, pass it through
							next.ServeHTTP(w, r.WithContext(ctx))
							return
						}
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func getAuthSignIn(w http.ResponseWriter, r *http.Request) {
	next := cleanNextURL(r.FormValue("next"))
	templates.Signin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Next":           next,
	})
}

func postAuthSignIn(w http.ResponseWriter, r *http.Request) {
	next := cleanNextURL(r.FormValue("next"))
	username := r.PostFormValue("username")
	passwd := r.PostFormValue("password")
	user := models.GetUserByPasswd(username, passwd)
	if user != nil {
		err := authenticate(user.ID, w)
		if err != nil {
			glog.Errorf("Error authenticating: %s", err.Error())
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
	}
	templates.Signin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Next":           next,
		"ErrMsg":         "Invalid username / password",
	})
}

func getAuthOneTimeSignIn(w http.ResponseWriter, r *http.Request) {
	next := cleanNextURL(r.FormValue("next"))
	templates.OneTimeSignin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Next":           next,
	})
}

func postAuthOneTimeSignIn(w http.ResponseWriter, r *http.Request) {
	next := cleanNextURL(r.PostFormValue("next"))
	email := r.PostFormValue("email")
	errMsg := "E-mail not found"

	users := models.GetUsersByEmail(email)
	for _, user := range *users {
		if len(user.Email) > 0 {
			errMsg = "A one time sign-in link has been sent to your email"
			token := models.UpdateUserOneTimeLoginTokenByID(user.ID)
			link := "http://" + r.Host + "/auth/otsignin/" + token + "?next=" + next

			forumName := models.GetConfigValue(models.ForumName)
			subject := forumName + " sign-in link"
			body := "Someone (hopefully you) requested a sign-in link for " + forumName + ".\r\n" +
				"If you want to sign-in, visit " + link + "\r\n\r\nIf not, just ignore this message."
			sendMail(user.Email, subject, body)
		}
	}

	templates.OneTimeSignin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Next":           next,
		"ErrMsg":         errMsg,
	})
}

func getAuthOneTimeSignInDone(w http.ResponseWriter, r *http.Request) {
	next := cleanNextURL(r.FormValue("next"))
	token := chi.URLParam(r, "token")
	user := models.GetUserByOneTimeToken(token)
	if user != nil {
		if err := authenticate(user.ID, w); err == nil {
			http.Redirect(w, r, next, http.StatusSeeOther)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func getAuthSignUp(w http.ResponseWriter, r *http.Request) {
	next := cleanNextURL(r.FormValue("next"))
	templates.Signup.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"Next":           next,
	})
}

func postAuthSignUp(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	passwd := r.PostFormValue("password")
	passwd2 := r.PostFormValue("password2")
	email := r.PostFormValue("email")
	next := cleanNextURL(r.FormValue("next"))
	newUserID := uuid.New()

	errMsg := ""
	if len(username) < 2 {
		errMsg = "Username should have at least 2 characters"
	} else if userNameReg.ReplaceAllString(username, "") != username {
		errMsg = "Username should not have non-alphanumeric characters"
	} else if len(passwd) < 6 {
		errMsg = "Password should have at least 6 characters"
	} else if passwd != passwd2 {
		errMsg = "Passwords do not match"
	}

	if errMsg == "" {
		email = strings.Trim(email, " ")
		if !strings.Contains(email, "@") {
			email = ""
		}
		err := models.CreateUser("", email, "", passwd)
		if err != nil {
			glog.Errorf("Error creating user: %s", err.Error())
			errMsg = "Error during signup."
		}
	}

	if errMsg == "" {
		glog.Infof("Created user: %s", username)
		err := authenticate(newUserID, w)
		if err != nil {
			glog.Errorf("Error authenticating: %s", err.Error())
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
	}

	templates.Signup.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"ErrMsg":         errMsg,
		"Next":           next,
	})
}

func getAuthChangePass(w http.ResponseWriter, r *http.Request) {
	templates.ChangePass.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func postAuthChangePass(w http.ResponseWriter, r *http.Request) {
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

	oldUser := models.GetUserByPasswd(user.Username, oldPasswd)
	if oldUser == nil || oldUser.ID != user.ID {
		errMsg = "Old password incorrect"
	}

	if errMsg == "" {
		err := models.UpdateUserPasswdByID(user.ID, passwd)
		if err == nil {
			errMsg = "Changed password successfully"
		} else {
			errMsg = "Error changing password"
			glog.Errorf("Error changing password: %s", err.Error())
		}
	}
	templates.ChangePass.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"ErrMsg":         errMsg,
	})
}

func getAuthLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: "", Path: "/", Expires: time.Now().Add(-300 * time.Hour), HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "csrftoken", Value: "", Path: "/", Expires: time.Now().Add(-300 * time.Hour)})
	if user, ok := r.Context().Value(CtxUserKey).(*models.User); ok {
		models.LogOutUserByID(user.ID)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func init() {
	userNameReg, _ = regexp.Compile("[^a-zA-Z0-9]+")
	nextURLReg, _ = regexp.Compile("[^a-zA-Z0-9-/]+")
}
