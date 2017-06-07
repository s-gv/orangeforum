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
	sess := models.OpenSession(w, r)
	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		passwd := r.PostFormValue("passwd")
		passwdConfirm := r.PostFormValue("confirm")
		email := r.PostFormValue("email")
		if r.PostFormValue("csrf") != sess.CSRFToken {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}
		if len(userName) == 0 {
			sess.SetFlashMsg("Username should not be blank.")
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
		if len(passwd) < 8 {
			sess.SetFlashMsg("Password should have at least 8 characters.")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		if passwd != passwdConfirm {
			sess.SetFlashMsg("Passwords don't match.")
			http.Redirect(w, r, "/signup", http.StatusSeeOther)
			return
		}
		models.CreateUser(userName, passwd, email)
		sess.Authenticate(userName, passwd)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	templates.Render(w, "signup.html", map[string]interface{}{
		"Msg": sess.FlashMsg(),
		"CSRF": sess.CSRFToken,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	sess := models.OpenSession(w, r)
	redirectURL := r.FormValue("next")
	if redirectURL == "" {
		redirectURL = "/"
	}
	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		passwd := r.PostFormValue("passwd")
		if r.PostFormValue("csrf") != sess.CSRFToken {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}
		if sess.Authenticate(userName, passwd) {
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		} else {
			sess.SetFlashMsg("Incorrect username/password")
			http.Redirect(w, r, "/login?next="+redirectURL, http.StatusSeeOther)
			return
		}
	}
	templates.Render(w, "login.html", map[string]interface{}{
		"CSRF": sess.CSRFToken,
		"Msg": sess.FlashMsg(),
		"next": redirectURL,
	})
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