package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)


var db *sql.DB

var stmts = make(map[string]*sql.Stmt)

type Row struct {
	*sql.Row
}

type Rows struct {
	*sql.Rows
}

func Init(driverName string, dataSourceName string) {
	mydb, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Panicf("[ERROR] Error opening DB: %s\n", err)
	}
	db = mydb
	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA synchronous = FULL;")
	db.Exec("PRAGMA foreign_keys = ON;")
}

func makeStmt(query string) *sql.Stmt {
	stmt, ok := stmts[query]
	if !ok {
		stmt, err := db.Prepare(query)
		if err != nil {
			log.Panicf("[ERROR] Error making stmt: %s. Err msg: %s\n", query, err)
		}
		stmts[query] = stmt
		return stmt
	}
	return stmt

}

func QueryRow(query string, args ...interface{}) *Row {
	stmt := makeStmt(query)
	return &Row{stmt.QueryRow(args...)}
}

func Query(query string, args ...interface{}) *Rows {
	stmt := makeStmt(query)
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Panicf("[ERROR] Error scanning rows: %s\n", err)
	}
	return &Rows{rows}
}

func (r *Row) Scan(args ...interface{}) error {
	err := r.Row.Scan(args...)
	switch {
	case err == sql.ErrNoRows:
		return err
	case err != nil:
		log.Panicf("[ERROR] Error scanning row: %s\n", err)
	}
	return nil
}

func (rs *Rows) Scan(args ...interface{}) error {
	err := rs.Rows.Scan(args...)
	switch {
	case err == sql.ErrNoRows:
		return err
	case err != nil:
		log.Panicf("[ERROR] Error scanning row: %s\n", err)
	}
	return nil
}

func Exec(query string, args ...interface{}) {
	stmt := makeStmt(query)
	if _, err := stmt.Exec(args...); err != nil {
		log.Panicf("[ERROR] Error executing %s. Err msg: %s\n", query, err)
	}
}

func Version() int {
	row := db.QueryRow(`SELECT val FROM configs WHERE key="version";`)
	sval := "0"
	if err := row.Scan(&sval); err == nil {
		if ival, err := strconv.Atoi(sval); err == nil {
			return ival
		}
	}
	return 0
}