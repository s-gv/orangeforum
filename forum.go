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
	dbDriverPtr := flag.String("dbdriver", "sqlite3", "Database driver name. Choose between sqlite3 and postgres (default: sqlite3)")
	dbSourceName := flag.String("dbsource", "orangeforum.db", "Database source name. For sqlite3, specify file path (default: orangeforum.db)")

	flag.Parse()

	models.Init(*dbDriverPtr, *dbSourceName)


	if(*shouldMigrate) {
		models.Migrate()
		return
	}

	http.HandleFunc("/", views.IndexHandler)

	port := ":9123"
	log.Println("[INFO] Starting orangeforum at port", port)
	http.ListenAndServe(port, nil)
}
