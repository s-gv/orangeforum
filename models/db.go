package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)
var db *sql.DB

func Init(driverName string, dataSourceName string) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	rows, err := db.Exec("SELECT * FROM config")

	if err != nil {
		log.Fatalln(err, "-- DB migration needed. Run with --migrate flag.\n")
	}

	rows = rows
}