package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init(driverName string, dataSourceName string) {
	mydb, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	db = mydb
}