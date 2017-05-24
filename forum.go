package main

import (
	"net/http"
	"io"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	println("Starting orange forum...")
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
