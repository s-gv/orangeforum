package views

import (
)
import (
	"net/http"
	"github.com/s-gv/orangeforum/templates"
	"log"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/utils"
	"strings"
	"errors"
	"html/template"
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

func ErrForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "403 Forbidden", http.StatusForbidden)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	if r.URL.Path != "/" {
		ErrNotFoundHandler(w, r)
		return
	}
	sess := models.OpenSession(w, r)
	flashMsg := sess.FlashMsg()
	name := ""
	if userName, err := sess.UserName(); err == nil {
		name = userName
	}
	templates.Render(w, "index.html", map[string]interface{}{
		"Title": "Orange Forum",
		"IsUserValid": sess.IsUserValid(),
		"UserName": name,
		"Karma": 812,
		"Msg": flashMsg,
	})
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

	redirectURL := r.FormValue("next")
	if redirectURL == "" {
		redirectURL = "/"
	}
	if sess.IsUserValid() {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		passwd := r.PostFormValue("passwd")
		passwdConfirm := r.PostFormValue("confirm")
		email := r.PostFormValue("email")
		if len(userName) == 0 {
			sess.SetFlashMsg("Username cannot be blank.")
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
		models.CreateUser(userName, passwd, email)
		sess.Authenticate(userName, passwd)
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
	templates.Render(w, "signup.html", map[string]interface{}{
		"Msg": sess.FlashMsg(),
		"CSRF": sess.CSRFToken,
		"next": template.URL(redirectURL),
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

	redirectURL := r.FormValue("next")
	if redirectURL == "" {
		redirectURL = "/"
	}
	if sess.IsUserValid() {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		passwd := r.PostFormValue("passwd")
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
		"next": template.URL(redirectURL),
	})
}

func ChangePasswdHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

	userName, err := sess.UserName()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == "POST" {
		passwd := r.PostFormValue("passwd")
		newPasswd := r.PostFormValue("newpass")
		newPasswdConfirm := r.PostFormValue("confirm")
		if !sess.Authenticate(userName, passwd) {
			sess.SetFlashMsg("Current password incorrect.")
			http.Redirect(w, r, "/changepass", http.StatusSeeOther)
			return
		}
		if err := validatePasswd(newPasswd, newPasswdConfirm); err != nil {
			sess.SetFlashMsg(err.Error())
			http.Redirect(w, r, "/changepass", http.StatusSeeOther)
			return
		}
		if err := models.UpdateUserPasswd(userName, newPasswd); err != nil {
			log.Fatalf("[ERROR] Error changing password: %s\n", err)
		}
		sess.SetFlashMsg("Password change successful.")
		http.Redirect(w, r, "/changepass", http.StatusSeeOther)
		return
	}
	templates.Render(w, "changepass.html", map[string]interface{}{
		"CSRF": sess.CSRFToken,
		"Msg": sess.FlashMsg(),
	})
}

func ForgotPasswdHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

	if r.Method == "POST" {
		userName := r.PostFormValue("username")
		if userName == "" || !models.ProbeUser(userName) {
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
		resetLink := "https://" + r.Host + "/resetpass?r=" + models.CreateResetToken(userName)
		sub := forumName + " Password Recovery"
		msg := "Someone (hopefully you) requested we reset your password at " + forumName + ".\r\n" +
			"If you want to change it, visit "+resetLink+"\r\n\r\nIf not, just ignore this message."
		utils.SendMail(email, sub, msg)
		sess.SetFlashMsg("Password reset link sent to your e-mail.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return

	}
	templates.Render(w, "forgotpass.html", map[string]interface{}{
		"CSRF": sess.CSRFToken,
		"Msg": sess.FlashMsg(),
	})
}

func validatePasswd(passwd string, passwdConfirm string) error {
	if len(passwd) < 8 {
		return errors.New("Password should have at least 8 characters.")
	}
	if passwd != passwdConfirm {
		return errors.New("Passwords don't match.")
	}
	return nil
}

func ResetPasswdHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}

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
		"CSRF": sess.CSRFToken,
		"Msg": sess.FlashMsg(),
	})
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	sess.SetFlashMsg("hi there")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	models.ClearSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}
	if !sess.IsUserValid() {
		http.Redirect(w, r, "/login?next=/creategroup", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		groupName := r.PostFormValue("name")
		//groupDesc := r.PostFormValue("desc")
		if groupName == "" {
			sess.SetFlashMsg("Group name is empty.")
			http.Redirect(w, r, "/creategroup", http.StatusSeeOther)
			return
		}
		hasSpecial := false
		for _, ch := range groupName {
			if (ch < 'A' || ch > 'Z') && (ch < 'a' || ch > 'z') && ch != '-' && (ch < '0' || ch > '9') {
				hasSpecial = true
			}
		}
		if hasSpecial {
			sess.SetFlashMsg("Username can contain only english alphabets, numbers, and hyphen.")
			http.Redirect(w, r, "/creategroup", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/g/"+groupName, http.StatusSeeOther)
		return
	}

	templates.Render(w, "creategroup.html", map[string]interface{}{
		"CSRF": sess.CSRFToken,
		"Msg": sess.FlashMsg(),
	})
}

func AdminIndexHandler(w http.ResponseWriter, r *http.Request) {
	defer ErrServerHandler(w, r)
	sess := models.OpenSession(w, r)
	if r.Method == "POST" && r.PostFormValue("csrf") != sess.CSRFToken {
		ErrForbiddenHandler(w, r)
		return
	}
	if !sess.IsUserSuperAdmin() {
		ErrForbiddenHandler(w, r)
		return
	}

	templates.Render(w, "adminindex.html", map[string]interface{}{
		"CSRF": sess.CSRFToken,
		"Msg": sess.FlashMsg(),
	})
}