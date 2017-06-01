package main

import (
	"net/http"
	"github.com/s-gv/orangeforum/views"
	"log"
	"flag"
	"github.com/s-gv/orangeforum/models/db"
)

func main() {
	dbFileName := flag.String("dbname", "orangeforum.db", "Database file path (default: orangeforum.db)")
	port := flag.String("port", "9123", "Port to listen on (default: 9123)")
	shouldMigrate := flag.Bool("migrate", false, "Migrate DB (default: false)")

	flag.Parse()

	if *shouldMigrate {
		err := db.Init("sqlite3", *dbFileName, true)
		if err != nil {
			log.Fatal("[ERROR] Migration failed. ", err)
		}
		log.Println("[INFO] DB migration successful.")
		return
	}

	err := db.Init("sqlite3", *dbFileName, false)
	if err != nil {
		log.Fatalln("[ERROR]", err)
	}


	http.HandleFunc("/", views.IndexHandler)
	http.HandleFunc("/signup", views.SignupHandler)
	http.HandleFunc("/login", views.LoginHandler)

	log.Println("[INFO] Starting orangeforum on port", *port)
	http.ListenAndServe(":" + *port, nil)
}
