package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"log"
	"time"
	"fmt"
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

func runMigrationZero() {
	db.Exec(`CREATE TABLE config(key TEXT, val TEXT);`)
	db.Exec(`CREATE UNIQUE INDEX key_index on config(key);`)

	db.Exec(`CREATE TABLE user(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		username TEXT NOT NULL,
		       		passwdhash TEXT,
		       		email TEXT,
		       		about TEXT,
		       		karma INTEGER,
		       		is_banned BOOLEAN,
		       		is_warned BOOLEAN,
		       		is_admin BOOLEAN,
		       		created_date INTEGER,
		       		updated_date INTEGER
	);`)

	db.Exec(`CREATE TABLE subforum(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		name TEXT,
		       		desc TEXT,
		       		created_date INTEGER,
		       		updated_date INTEGER
	);`)

	db.Exec(`CREATE TABLE mod(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		user_id INTEGER,
		       		subforum_id INTEGER,
		       		created_date INTEGER
	);`)

	WriteConfig("version", "1")
}


func DBVersion() int {
	if val, err := strconv.Atoi(Config("version", "0")); err == nil {
		return val
	}
	return 0

}

func WriteConfig(key string, val string) error {
	return exec("WriteConfig", `INSERT OR REPLACE INTO config(key, val) values(?, ?);`, key, val)
}


func Config(key string, defaultVal string) string {
	row, err := queryRow("Config", `SELECT val FROM config WHERE key=?;`, "version")
	if err == nil {
		var val string
		if err := row.Scan(&val); err == nil {
			return val
		}
	}
	return defaultVal
}

func Init(driverName string, dataSourceName string) error {
	mydb, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	db = mydb
	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA synchronous = FULL;")

	dbver := DBVersion()
	if dbver < ModelVersion {
		return ErrDBVer
	}
	return nil
}

func Benchmark() {
	start := time.Now()
	for i := 0; i < 1000; i++ {
		WriteConfig("version", "3")
	}
	elapsed := time.Since(start)
	WriteConfig("version", "2")
	println("Test val:", Config("version", "0"))
	fmt.Printf("time: %s\n", elapsed)
}