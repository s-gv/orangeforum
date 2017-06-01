package views

import (
	"net/http"
	"github.com/s-gv/orangeforum/templates"
	"github.com/s-gv/orangeforum/models/sessions"
	"github.com/s-gv/orangeforum/models"
	"log"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sess := sessions.Open(w, r)
	flashMsg := sess.FlashMsg()
	name := "world"
	if u, err := sess.User(); err == nil {
		name = u.Username
	} else {
		log.Println(err)
	}
	templates.Render(w, "index.html", map[string]interface{}{
		"Title": "Orange Forum",
		"Name": name,
		"Msg": flashMsg,
	})
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	sess := sessions.Open(w, r)
	if r.Method == "POST" {
		userName := "deaf"
		passwd := "1234"
		email := "fda@gdafdas.com"
		if models.ProbeUser(userName) {
			sess.SetFlashMsg("User already exists.")
		} else {
			models.CreateUser(userName, passwd, email)
			sess.SetFlashMsg("Created user successfully")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sess := sessions.Open(w, r)
	userName := "deaf"
	passwd := "1234"
	if user, err := models.Authenticate(userName, passwd); err == nil {
		sess.SetUser(user, true)
		sess.SetFlashMsg("Logged in")
	} else {
		sess.SetFlashMsg(err.Error())
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}