package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

var db *sql.DB

func DBVersion() int {
	if val, err := strconv.Atoi(Config("version", "0")); err == nil {
		return val
	}
	return 0

}

func createConfigTable() {
	db.Exec(`CREATE TABLE config(key TEXT NOT NULL, val TEXT)`)

}

func WriteConfig(key string, val string) {
	db.Exec(`INSERT INTO config(key, val) values(?, ?)`, key, val)
}

func Config(key string, defaultVal string) string {
	row := db.QueryRow(`SELECT val FROM config WHERE key=?`, "version")
	var val string
	if err := row.Scan(&val); err == nil {
		return val
	}
	return defaultVal
}