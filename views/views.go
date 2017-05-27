package views

import (
	"net/http"
	"github.com/s-gv/orangeforum/templates"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	templates.Render(w, "index.html", map[string]interface{}{
		"Title": "Orange Forum",
	})
}