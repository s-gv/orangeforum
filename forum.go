package main

import (
	"net/http"
	"github.com/s-gv/orangeforum/views"
	"log"
	"github.com/s-gv/orangeforum/models"
	"flag"
	"os"
)

func main() {
	dbFileName := flag.String("dbname", "orangeforum.db", "Database file path (default: orangeforum.db)")

	flag.Parse()

	if len(os.Args) > 1 {
		if os.Args[1] == "migrate" {
			err := models.Init("sqlite3", *dbFileName, true)
			if err != nil {
				log.Fatal("[ERROR] DB migration failed. ", err)
			}
			log.Println("[INFO] DB migration successful.")
			return
		}
	}

	err := models.Init("sqlite3", *dbFileName, false)
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}


	http.HandleFunc("/", views.IndexHandler)
	http.HandleFunc("/test", views.TestHandler)

	port := ":9123"
	log.Println("[INFO] Starting orangeforum at port", port)
	http.ListenAndServe(port, nil)
}
