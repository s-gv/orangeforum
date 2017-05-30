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

	err := models.Init("sqlite3", *dbFileName)

	if len(os.Args) > 1 {
		if os.Args[1] == "migrate" {
			err := models.Migrate()
			if err != nil {
				log.Fatal("[ERROR] ", err)

			}
			log.Println("[INFO] DB migration successful.")
			return
		}
		if os.Args[1] == "benchmark" {
			models.Benchmark()
			return
		}

	}

	if err != nil {
		log.Fatal("[ERROR] ", err)
	}


	http.HandleFunc("/", views.IndexHandler)

	port := ":9123"
	log.Println("[INFO] Starting orangeforum at port", port)
	http.ListenAndServe(port, nil)
}
