package main

import (
	"net/http"
	"github.com/s-gv/orangeforum/templates"
)

func hello(w http.ResponseWriter, r *http.Request) {
	templates.Render(w, "index.html", map[string]interface{}{
		"Title": "Orange Forum",
	})
}

func main() {
	println("Starting orange forum...")
	http.HandleFunc("/", hello)
	http.ListenAndServe(":9123", nil)
}
