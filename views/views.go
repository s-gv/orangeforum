package views

import (
	"net/http"
	"github.com/s-gv/orangeforum/templates"
	"github.com/s-gv/orangeforum/models"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	flashMsg := ""
	if msg, err := models.GetFlashMsg(w, r); err == nil {
		flashMsg = msg
	}
	templates.Render(w, "index.html", map[string]interface{}{
		"Title": "Orange Forum",
		"Msg": flashMsg,
	})
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		models.SetFlashMsg(w, "This is a flash message.")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}