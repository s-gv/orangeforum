package main

import (
	"net/http"
	"github.com/s-gv/orangeforum/views"
	"log"
)

func main() {
	http.HandleFunc("/", views.IndexHandler)

	port := ":9123"
	log.Println("[INFO] Starting orangeforum at port", port)
	http.ListenAndServe(port, nil)
}
