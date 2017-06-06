package views

import (
)
import (
	"net/http"
	"github.com/s-gv/orangeforum/templates"
	"log"
	"github.com/s-gv/orangeforum/models"
)

func ErrServerHandler(w http.ResponseWriter, r *http.Request) {
	if r := recover(); r != nil {
		log.Printf("[INFO] Recovered from panic: %s", r)
		http.Error(w, "Internal server error. This event has been logged.", http.StatusInternalServerError)
	}
}

func ErrNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	if r.URL.Path != "/" {
		ErrNotFoundHandler(w, r)
		return
	}
	sess := models.OpenSession(w, r)
	flashMsg := sess.FlashMsg()
	name := "world"
	if userName, err := sess.UserName(); err == nil {
		name = userName
	}
	templates.Render(w, "index.html", map[string]interface{}{
		"Title": "Orange Forum",
		"Name": name,
		"Msg": flashMsg,
	})
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sess := models.OpenSession(w, r)
	userName := "deaf"
	passwd := "1234"
	email := "fda@gdafdas.com"
	if err := models.CreateUser(userName, passwd, email); err == nil {
		sess.SetFlashMsg("Created user successfully")
	} else {
		sess.SetFlashMsg(err.Error())
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sess := models.OpenSession(w, r)
	userName := "deaf"
	passwd := "1234"
	if sess.Authenticate(userName, passwd) {
		sess.SetFlashMsg("Logged in")
	} else {
		sess.SetFlashMsg("Incorrect username/password")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	sess := models.OpenSession(w, r)
	sess.SetFlashMsg("hi there")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	models.ClearSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}