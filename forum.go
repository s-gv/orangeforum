package main

import (
	"net/http"
	"log"
	"flag"
	"github.com/s-gv/orangeforum/models/db"
	"time"
	"math/rand"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/views"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	dbFileName := flag.String("dbname", "orangeforum.db", "Database file path (default: orangeforum.db)")
	port := flag.String("port", "9123", "Port to listen on (default: 9123)")
	shouldMigrate := flag.Bool("migrate", false, "Migrate DB (default: false)")

	flag.Parse()

	db.Init("sqlite3", *dbFileName)

	if models.IsMigrationNeeded() {
		if *shouldMigrate {
			models.Migrate()
		} else {
			log.Fatalf("[ERROR] DB migration needed.\n")
		}
	} else {
		if *shouldMigrate {
			log.Fatalf("[ERROR] DB migration not needed. DB up-to-date.\n")
		}
	}


	http.HandleFunc("/", views.IndexHandler)
	/*
	http.HandleFunc("/signup", views.SignupHandler)
	http.HandleFunc("/login", views.LoginHandler)
	*/

	log.Println("[INFO] Starting orangeforum on port", *port)
	http.ListenAndServe(":" + *port, nil)
}
