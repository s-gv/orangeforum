// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"fmt"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/models/db"
	"github.com/s-gv/orangeforum/templates"
	"github.com/s-gv/orangeforum/utils"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var LoginHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	redirectURL, err := url.QueryUnescape(r.FormValue("next"))
	if err != nil || redirectURL == "" || redirectURL[0] != '/' {
		redirectURL = "/"
	}
	if sess.IsUserValid() {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		passwd := r.PostFormValue("passwd")
		if len(userName) > 200 || len(passwd) > 200 {
			fmt.Fprint(w, "username / password too long.")
			return
		}
		if err = sess.Authenticate(userName, passwd); err == nil {
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		} else {
			sess.SetFlashMsg(err.Error())
			http.Redirect(w, r, "/login?next="+redirectURL, http.StatusSeeOther)
			return
		}
	}
	templates.Render(w, "login.html", map[string]interface{}{
		"Common":   readCommonData(r, sess),
		"next":     template.URL(url.QueryEscape(redirectURL)),
		"LoginMsg": models.Config(models.LoginMsg),
	})
})

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	ClearSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

var SignupHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	redirectURL, err := url.QueryUnescape(r.FormValue("next"))
	if err != nil || redirectURL == "" || redirectURL[0] != '/' {
		redirectURL = "/"
	}
	if sess.IsUserValid() && !sess.IsUserSuperAdmin() {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	isSignupDisabled := models.Config(models.SignupDisabled) != "0"

	if r.Method == "POST" {
		userName := strings.TrimSpace(r.PostFormValue("username"))
		passwd := r.PostFormValue("passwd")
		passwdConfirm := r.PostFormValue("confirm")
		email := strings.TrimSpace(r.PostFormValue("email"))
		if len(userName) < 2 || len(userName) > 32 {
			sess.SetFlashMsg("Username should have 2-32 characters.")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		if censored := censor(userName); censored != userName {
			sess.SetFlashMsg("Fix username: " + censored)
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		hasSpecial := false
		for _, ch := range userName {
			if (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') && ch != '_' && (ch < '0' || ch > '9') {
				hasSpecial = true
			}
		}
		if hasSpecial {
			sess.SetFlashMsg("Username can contain only alphabets, numbers, and underscore.")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		if models.ProbeUser(userName) {
			sess.SetFlashMsg("Username already registered.")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		if err := validatePasswd(passwd, passwdConfirm); err != nil {
			sess.SetFlashMsg(err.Error())
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		if len(email) > 64 {
			sess.SetFlashMsg("Email should have fewer than 64 characters.")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		if isSignupDisabled && !sess.IsUserSuperAdmin() {
			ErrForbiddenHandler(w, r)
			return
		}
		models.CreateUser(userName, passwd, email)
		if sess.IsUserSuperAdmin() {
			sess.SetFlashMsg("User " + userName + " created")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		sess.Authenticate(userName, passwd)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
	templates.Render(w, "signup.html", map[string]interface{}{
		"Common":     readCommonData(r, sess),
		"next":       template.URL(url.QueryEscape(redirectURL)),
		"IsDisabled": isSignupDisabled && !sess.IsUserSuperAdmin(),
		"SignupMsg":  models.Config(models.SignupMsg),
	})
})

var ChangePasswdHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	userName := r.FormValue("u")
	commonData := readCommonData(r, sess)
	if !sess.IsUserValid() {
		ErrForbiddenHandler(w, r)
		return
	}
	if userName != commonData.UserName && !commonData.IsSuperAdmin {
		ErrForbiddenHandler(w, r)
		return
	}
	if r.Method == "POST" {
		if !commonData.IsSuperAdmin {
			passwd := r.PostFormValue("passwd")
			if sess.Authenticate(userName, passwd) != nil {
				sess.SetFlashMsg("Current password incorrect.")
				http.Redirect(w, r, "/changepass?u="+userName, http.StatusSeeOther)
				return
			}
		}
		newPasswd := r.PostFormValue("newpass")
		newPasswdConfirm := r.PostFormValue("confirm")
		if err := validatePasswd(newPasswd, newPasswdConfirm); err != nil {
			sess.SetFlashMsg(err.Error())
			http.Redirect(w, r, "/changepass?u="+userName, http.StatusSeeOther)
			return
		}
		if err := models.UpdateUserPasswd(userName, newPasswd); err != nil {
			log.Panicf("[ERROR] Error changing password: %s\n", err)
		}
		if commonData.IsSuperAdmin {
			var userID string
			db.QueryRow(`SELECT id FROM users WHERE username=?;`, userName).Scan(&userID)
			db.Exec(`DELETE FROM sessions WHERE userid=?;`, userID)
		}
		sess.SetFlashMsg("Password change successful.")
		http.Redirect(w, r, "/changepass?u="+userName, http.StatusSeeOther)
		return
	}
	templates.Render(w, "changepass.html", map[string]interface{}{
		"Common":   commonData,
		"UserName": userName,
	})
})

var ForgotPasswdHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		if userName == "" || len(userName) > 200 || !models.ProbeUser(userName) {
			sess.SetFlashMsg("Username doesn't exist.")
			http.Redirect(w, r, "/forgotpass", http.StatusSeeOther)
			return
		}
		email := models.ReadUserEmail(userName)
		if !strings.ContainsRune(email, '@') {
			sess.SetFlashMsg("E-mail address not set. Contact site admin to reset the password.")
			http.Redirect(w, r, "/forgotpass", http.StatusSeeOther)
			return
		}
		forumName := models.Config(models.ForumName)

		resetToken := randSeq(40)
		db.Exec(`UPDATE users SET reset_token=?, reset_token_date=? WHERE username=?;`, resetToken, int64(time.Now().Unix()), userName)

		resetLink := "https://" + r.Host + "/resetpass?r=" + resetToken
		sub := forumName + " Password Recovery"
		msg := "Someone (hopefully you) requested we reset your password at " + forumName + ".\r\n" +
			"If you want to change it, visit " + resetLink + "\r\n\r\nIf not, just ignore this message."
		utils.SendMail(email, sub, msg)
		sess.SetFlashMsg("Password reset link sent to your e-mail.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	}
	templates.Render(w, "forgotpass.html", map[string]interface{}{
		"Common": readCommonData(r, sess),
	})
})

var ResetPasswdHandler = UA(func(w http.ResponseWriter, r *http.Request, sess Session) {
	resetToken := r.FormValue("r")
	userName, err := models.ReadUserNameByToken(resetToken)
	if err != nil {
		ErrForbiddenHandler(w, r)
		return
	}
	if r.Method == "POST" {
		passwd := r.PostFormValue("passwd")
		passwdConfirm := r.PostFormValue("confirm")
		if err := validatePasswd(passwd, passwdConfirm); err != nil {
			sess.SetFlashMsg(err.Error())
			http.Redirect(w, r, "/resetpass?r="+resetToken, http.StatusSeeOther)
			return
		}
		models.UpdateUserPasswd(userName, passwd)
		sess.SetFlashMsg("Password change successful.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.Render(w, "resetpass.html", map[string]interface{}{
		"ResetToken": resetToken,
		"Common":     readCommonData(r, sess),
	})
})
