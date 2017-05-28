package main

import (
	"net/http"
	"github.com/s-gv/orangeforum/views"
	"log"
	"github.com/s-gv/orangeforum/models"
	"flag"
)

func main() {
	shouldMigrate := flag.Bool("migrate", false, "Migrate DB to the current version (default: false)")
	benchmark := flag.Bool("benchmark", false, "Run the benchmark")
	dbFileName := flag.String("dbname", "orangeforum.db", "Database file path (default: orangeforum.db)")

	flag.Parse()

	err := models.Init("sqlite3", *dbFileName)
	if *shouldMigrate {
		err := models.Migrate()
		if err != nil {
			log.Fatal("[ERROR] ", err)

		}
		log.Println("[INFO] DB migration successful.")
		return
	}
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}

	if *benchmark {
		models.Benchmark()
		return
	}

	http.HandleFunc("/", views.IndexHandler)

	port := ":9123"
	log.Println("[INFO] Starting orangeforum at port", port)
	http.ListenAndServe(port, nil)
}
