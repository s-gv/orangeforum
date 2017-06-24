package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"log"
	"strconv"
	"strings"
)


var db *sql.DB
var dbDriverName string

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
	dbDriverName = driverName
	if driverName == "sqlite3" {
		db.Exec("PRAGMA journal_mode = WAL;")
		db.Exec("PRAGMA synchronous = FULL;")
		db.Exec("PRAGMA foreign_keys = ON;")
	}
}

func translate(query string) string {
	if dbDriverName == "postgres" {
		query = strings.Replace(query, "INTEGER PRIMARY KEY AUTOINCREMENT", "SERIAL PRIMARY KEY", -1)
		query = strings.Replace(query, "MAX(_ROWID_)", "COUNT(*)", -1)
		p := 0
		for i := strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
			p++
			query = query[:i] + "$" + strconv.Itoa(p) + query[i+1:]
		}
	}
	return query
}

func patch(args []interface{}) []interface{} {
	var pArgs []interface{}
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			if v {
				pArgs = append(pArgs, 1)
			} else {
				pArgs = append(pArgs, 0)
			}
		default:
			pArgs = append(pArgs, v)
		}
	}
	return pArgs
}

func makeStmt(query string) *sql.Stmt {
	stmt, ok := stmts[query]
	if !ok {
		stmt, err := db.Prepare(translate(query))
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
	pArgs := patch(args)
	return &Row{stmt.QueryRow(pArgs...)}
}

func Query(query string, args ...interface{}) *Rows {
	stmt := makeStmt(query)
	pArgs := patch(args)
	rows, err := stmt.Query(pArgs...)
	if err != nil {
		log.Panicf("[ERROR] Error with SQL query '%s': %s\n", query, err)
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
		log.Panicf("[ERROR] Error scanning rows: %s\n", err)
	}
	return nil
}

func Exec(query string, args ...interface{}) {
	pArgs := patch(args)
	stmt := makeStmt(query)
	if _, err := stmt.Exec(pArgs...); err != nil {
		log.Panicf("[ERROR] Error executing %s. Err msg: %s\n", query, err)
	}
}

func Version() int {
	row := db.QueryRow(translate(`SELECT val FROM configs WHERE name=?;`), "version")
	sval := "0"
	if err := row.Scan(&sval); err == nil {
		if ival, err := strconv.Atoi(sval); err == nil {
			return ival
		}
	}
	return 0
}