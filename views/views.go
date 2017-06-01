package views

import (
	"net/http"
	"github.com/s-gv/orangeforum/templates"
	"github.com/s-gv/orangeforum/models/sessions"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	sess := sessions.Open(w, r)
	flashMsg := sess.FlashMsg()
	templates.Render(w, "index.html", map[string]interface{}{
		"Title": "Orange Forum",
		"Msg": flashMsg,
	})
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	sess := sessions.Open(w, r)
	if r.Method == "POST" {
		sess.SetFlashMsg("This is a flash message.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}