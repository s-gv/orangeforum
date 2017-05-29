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
	if _, err := db.Exec(`CREATE TABLE config(key TEXT, val TEXT);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE UNIQUE INDEX key_index on config(key);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE user(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		username TEXT NOT NULL,
		       		passwdhash TEXT,
		       		email TEXT,
		       		about TEXT,
		       		karma INTEGER,
		       		is_banned BOOLEAN,
		       		is_warned BOOLEAN,
		       		is_admin BOOLEAN,
				is_supermod BOOLEAN,
		       		created_date INTEGER,
		       		updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE UNIQUE INDEX username_index on user(username);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX email_index on user(email);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE category(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		name TEXT,
		       		desc TEXT,
		       		created_date INTEGER,
		       		updated_date INTEGER
	);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE mod(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		userid INTEGER REFERENCES user(id) ON DELETE CASCADE,
				categoryid INTEGER REFERENCES category(id) ON DELETE CASCADE,
		       		created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX mod_userid_index on mod(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX mod_categoryid_index on mod(categoryid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE topic(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT,
				authorid INTEGER REFERENCES user(id) ON DELETE CASCADE,
				categoryid INTEGER REFERENCES category(id) ON DELETE SET NULL,
				is_deleted BOOLEAN,
				is_closed BOOLEAN,
				is_sticky BOOLEAN,
				upvotes INTEGER,
				downvotes INTEGER,
				flagvotes INTEGER,
				created_date INTEGER,
				updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topic_authorid_index on topic(authorid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topic_categoryid_index on topic(categoryid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE comment(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT,
				authorid INTEGER REFERENCES user(id) ON DELETE CASCADE,
				topicid INTEGER REFERENCES topic(id) ON DELETE CASCADE,
				parentid INTEGER REFERENCES comment(id) ON DELETE CASCADE,
				is_deleted TEXT,
				upvotes INTEGER,
				downvotes INTEGER,
				flagvotes INTEGER,
				created_date INTEGER,
				updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX comment_authorid_index on comment(authorid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX comment_topicid_index on comment(topicid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX comment_parentid_index on comment(parentid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE topicvote(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				userid INTEGER REFERENCES user(id) ON DELETE CASCADE,
				topicid INTEGER REFERENCES topic(id) ON DELETE CASCADE,
				votetype INTEGER,
				created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topicvote_userid_index on topicvote(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topicvote_topicid_index on topicvote(topicid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE commentvote(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				userid INTEGER REFERENCES user(id) ON DELETE CASCADE,
				commentid INTEGER REFERENCES comment(id) ON DELETE CASCADE,
				votetype INTEGER,
				created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX commentvote_userid_index on commentvote(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX commentvote_commentid_index on commentvote(commentid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE session(
				id INTEGER PRIMARY KEY,
				sessionid TEXT NOT NULL,
				userid INTEGER REFERENCES user(id) ON DELETE CASCADE,
				data TEXT,
				created_date INTEGER,
				updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE UNIQUE INDEX sessionid_index on session(sessionid);`); err != nil { panic(err) }

	
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
	db.Exec("PRAGMA foreign_keys = ON;")

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