package models

import (
	"errors"
	"net/http"
	"time"
)

var ErrAuthFail = errors.New("username / password incorrect")
var ErrNoFlashMsg = errors.New("No flash message")

func Authenticate() error {
	return nil
}

func SetFlashMsg(w http.ResponseWriter, msg string) {
	cookie := http.Cookie{Name: "flashMsg", Value: msg, HttpOnly: true}
	http.SetCookie(w, &cookie)
}

func GetFlashMsg(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("flashMsg")
	http.SetCookie(w, &http.Cookie{Name: "flashMsg", Value: "", Expires: time.Now().Add(-2000*time.Hour)})
	if err == nil {
		return cookie.Value, nil
	}
	return "", ErrNoFlashMsg
}