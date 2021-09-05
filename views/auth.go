// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"context"
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

func cleanNextURL(next string, basePath string) string {
	if next == "" || next[0] != '/' {
		return basePath
	}
	return next
}

func authenticate(id int, basePath string, w http.ResponseWriter) error {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"user_id": strconv.Itoa(id),
		"iat":     time.Now(),
		"exp":     time.Now().Add(365 * 24 * time.Hour),
	})
	path := "/"
	if basePath != "" {
		path = basePath
	}
	if err == nil {
		cookie := http.Cookie{
			Name:     "jwt",
			Value:    tokenString,
			Path:     path,
			Expires:  time.Now().Add(365 * 24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
	}
	return err
}

func mustAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainID := r.Context().Value(ctxDomain).(*models.Domain).DomainID
		token, claims, err := jwtauth.FromContext(r.Context())
		basePath, _ := r.Context().Value(ctxBasePath).(string)

		if err == nil && token != nil && jwt.Validate(token) == nil {
			if uid, ok := claims["user_id"].(string); ok {
				if iat, ok := claims["iat"].(time.Time); ok {
					userID, _ := strconv.Atoi(uid)
					user := models.GetUserByID(userID)
					if user != nil && user.LogoutAt.Before(iat) && user.DomainID == domainID && !user.BannedAt.Valid {
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
			http.Redirect(w, r, basePath+"auth/signin?next="+r.URL.Path+"?"+r.URL.RawQuery, http.StatusSeeOther)
		} else {
			http.Redirect(w, r, basePath, http.StatusSeeOther)
		}
	})
}

func canAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainID := r.Context().Value(ctxDomain).(*models.Domain).DomainID
		token, claims, err := jwtauth.FromContext(r.Context())

		if err == nil && token != nil && jwt.Validate(token) == nil {
			if uid, ok := claims["user_id"].(string); ok {
				if iat, ok := claims["iat"].(time.Time); ok {
					userID, _ := strconv.Atoi(uid)
					user := models.GetUserByID(userID)
					if user != nil && user.LogoutAt.Before(iat) && user.DomainID == domainID {
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
	basePath := r.Context().Value(ctxBasePath).(string)
	domain := r.Context().Value(ctxDomain).(*models.Domain)
	next := cleanNextURL(r.FormValue("next"), basePath)
	templates.Signin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		DomainField:      domain,
		"Next":           next,
	})
}

func postAuthSignIn(w http.ResponseWriter, r *http.Request) {
	domain := r.Context().Value(ctxDomain).(*models.Domain)
	basePath := r.Context().Value(ctxBasePath).(string)
	next := cleanNextURL(r.FormValue("next"), basePath)
	email := r.PostFormValue("email")
	passwd := r.PostFormValue("password")
	user := models.GetUserByPasswd(domain.DomainID, email, passwd)

	if user != nil && !user.BannedAt.Valid {
		err := authenticate(user.UserID, basePath, w)
		if err != nil {
			glog.Errorf("Error authenticating: %s", err.Error())
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
	}
	templates.Signin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		DomainField:      domain,
		"Next":           next,
		"ErrMsg":         "Invalid email / password",
	})
}

func getAuthOneTimeSignIn(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)
	domain := r.Context().Value(ctxDomain).(*models.Domain)
	next := cleanNextURL(r.FormValue("next"), basePath)
	templates.OneTimeSignin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		DomainField:      domain,
		"Next":           next,
	})
}

func postAuthOneTimeSignIn(w http.ResponseWriter, r *http.Request) {
	domain := r.Context().Value(ctxDomain).(*models.Domain)
	basePath := r.Context().Value(ctxBasePath).(string)
	next := cleanNextURL(r.PostFormValue("next"), basePath)
	email := r.PostFormValue("email")
	errMsg := "E-mail not found"

	user := models.GetUserByEmail(domain.DomainID, email)
	if user != nil {
		errMsg = "A one time sign-in link has been sent to your email"
		token := models.UpdateUserOneTimeLoginTokenByID(user.UserID)
		link := "http://" + r.Host + basePath + "auth/otsignin/" + token + "?next=" + next

		forumName := domain.ForumName
		subject := forumName + " sign-in link"
		body := "Someone (hopefully you) requested a sign-in link for " + forumName + ".\r\n" +
			"If you want to sign-in, visit " + link + "\r\n\r\nIf not, just ignore this message."
		sendMail(domain.DefaultFromEmail, user.Email, subject, body, domain.ForumName, domain.SMTPHost, domain.SMTPPort, domain.SMTPUser, domain.SMTPPass)
	}

	templates.OneTimeSignin.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		DomainField:      domain,
		"Next":           next,
		"ErrMsg":         errMsg,
	})
}

func getAuthOneTimeSignInDone(w http.ResponseWriter, r *http.Request) {
	domainID := r.Context().Value(ctxDomain).(*models.Domain).DomainID
	basePath := r.Context().Value(ctxBasePath).(string)
	next := cleanNextURL(r.FormValue("next"), basePath)
	token := chi.URLParam(r, "token")
	user := models.GetUserByOneTimeToken(domainID, token)
	if user != nil {
		if err := authenticate(user.UserID, basePath, w); err == nil {
			http.Redirect(w, r, next, http.StatusSeeOther)
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func getAuthSignUp(w http.ResponseWriter, r *http.Request) {
	domain := r.Context().Value(ctxDomain).(*models.Domain)
	basePath := r.Context().Value(ctxBasePath).(string)
	next := cleanNextURL(r.FormValue("next"), basePath)
	signupToken := chi.URLParam(r, "signupToken")

	if signupToken != "" && signupToken != domain.SignupToken {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}

	if !domain.IsRegularSignupEnabled {
		if domain.SignupToken == "" || signupToken == "" || domain.SignupToken != signupToken {
			http.Error(w, "Signup is disabled. Contact the admin.", http.StatusForbidden)
			return
		}
	}
	templates.Signup.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		DomainField:      domain,
		"Next":           next,
	})
}

func validateEmail(email string) string {
	errMsg := ""
	if !strings.Contains(email, "@") {
		errMsg = "Invalid email"
	}
	if len(email) < 3 {
		errMsg = "Email should have at least 3 characters"
	} else if emailReg.ReplaceAllString(email, "") != email {
		errMsg = "Email should not have non-alphanumeric characters"
	}
	return errMsg
}

func postAuthSignUp(w http.ResponseWriter, r *http.Request) {
	domain := r.Context().Value(ctxDomain).(*models.Domain)
	basePath := r.Context().Value(ctxBasePath).(string)
	next := cleanNextURL(r.FormValue("next"), basePath)
	signupToken := chi.URLParam(r, "signupToken")

	email := r.PostFormValue("email")
	passwd := r.PostFormValue("password")
	passwd2 := r.PostFormValue("password2")

	email = strings.Trim(email, " ")

	if signupToken != "" && signupToken != domain.SignupToken {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}

	errMsg := ""

	errMsg = validateEmail(email)

	if len(passwd) < 6 {
		errMsg = "Password should have at least 6 characters"
	} else if passwd != passwd2 {
		errMsg = "Passwords do not match"
	}

	existingUser := models.GetUserByEmail(domain.DomainID, email)
	if existingUser != nil {
		errMsg = "E-mail already registered"
	}

	if !domain.IsRegularSignupEnabled {
		if domain.SignupToken == "" || signupToken == "" || domain.SignupToken != signupToken {
			http.Error(w, "Signup is disabled. Contact the admin.", http.StatusForbidden)
			return
		}
	}

	if errMsg == "" {
		displayName := strings.Split(email, "@")[0]
		displayName = strings.Title(strings.ToLower(displayName))

		err := models.CreateUser(domain.DomainID, email, displayName, passwd)
		if err != nil {
			glog.Errorf("Error creating user: %s", err.Error())
			errMsg = "Error during signup."
		}
	}

	if errMsg == "" {
		glog.Infof("Created user: %s for domainID: %d", email, domain.DomainID)
		user := models.GetUserByEmail(domain.DomainID, email)
		err := authenticate(user.UserID, basePath, w)
		if err != nil {
			glog.Errorf("Error authenticating: %s", err.Error())
		}
		http.Redirect(w, r, next, http.StatusSeeOther)
	}

	templates.Signup.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		DomainField:      domain,
		"ErrMsg":         errMsg,
		"Next":           next,
	})
}

func getAuthChangePass(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)
	domain := r.Context().Value(ctxDomain).(*models.Domain)
	user := r.Context().Value(CtxUserKey).(*models.User)
	templates.ChangePass.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath,
		DomainField:      domain,
		UserField:        user,
	})
}

func postAuthChangePass(w http.ResponseWriter, r *http.Request) {
	domain := r.Context().Value(ctxDomain).(*models.Domain)
	basePath := r.Context().Value(ctxBasePath).(string)

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

	oldUser := models.GetUserByPasswd(domain.DomainID, user.Email, oldPasswd)
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
		BasePathField:    basePath,
		DomainField:      domain,
		UserField:        user,
		"ErrMsg":         errMsg,
	})
}

func getAuthLogout(w http.ResponseWriter, r *http.Request) {
	basePath := r.Context().Value(ctxBasePath).(string)

	http.SetCookie(w, &http.Cookie{Name: "jwt", Value: "", Path: basePath, Expires: time.Now().Add(-300 * time.Hour), HttpOnly: true})
	http.SetCookie(w, &http.Cookie{Name: "csrftoken", Value: "", Path: basePath, Expires: time.Now().Add(-300 * time.Hour)})
	if user, ok := r.Context().Value(CtxUserKey).(*models.User); ok {
		models.LogOutUserByID(user.UserID)
	}
	http.Redirect(w, r, basePath, http.StatusSeeOther)
}

func init() {
	userNameReg = regexp.MustCompile("[^a-zA-Z0-9]+")
	emailReg = regexp.MustCompile("[^a-zA-Z0-9@\\.-]+")
	nextURLReg = regexp.MustCompile("[^a-zA-Z0-9-/]+")
}
