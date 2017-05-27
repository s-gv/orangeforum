package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"log"
)

var db *sql.DB

var stmts = make(map[string]*sql.Stmt)

func makeStmt(name string, query string) (*sql.Stmt, error) {
	stmt, ok := stmts[name]
	if !ok {
		var err error
		stmt, err := db.Prepare(query)
		if err != nil {
			return nil, err
		}
		stmts[name] = stmt
		return stmt, nil
	}
	return stmt, nil

}

func queryRow(name string, query string, args ...interface{}) (*sql.Row, error) {
	stmt, err := makeStmt(name, query)
	if err == nil {
		return stmt.QueryRow(args...), nil
	}
	return nil, err
}

func exec(name string, query string, args ...interface{}) error {
	stmt, err := makeStmt(name, query)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(args...); err != nil {
		log.Println("[ERROR] Error executing", name, err)
		return err
	}
	return nil
}

func createConfigTable() {
	db.Exec(`CREATE TABLE config(key TEXT NOT NULL, val TEXT)`)

}

func DBVersion() int {
	if val, err := strconv.Atoi(Config("version", "0")); err == nil {
		return val
	}
	return 0

}


func WriteConfig(key string, val string) {
	exec("WriteConfig", `INSERT INTO config(key, val) values(?, ?)`, key, val)
}


func Config(key string, defaultVal string) string {
	row, err := queryRow("Config", `SELECT val FROM config WHERE key=?`, "version")
	if err == nil {
		var val string
		if err := row.Scan(&val); err == nil {
			return val
		}
	}
	return defaultVal
}