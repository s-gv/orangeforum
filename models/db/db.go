package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"log"
)


var db *sql.DB

var stmts = make(map[string]*sql.Stmt)

type Row struct {
	*sql.Row
}

type Rows struct {
	*sql.Rows
}

func CreateTables() {
	if _, err := db.Exec(`CREATE TABLE configs(key VARCHAR(250), val TEXT);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE UNIQUE INDEX configs_key_index on configs(key);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE users(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		username VARCHAR(32) NOT NULL,
		       		passwdhash VARCHAR(250) NOT NULL,
		       		email VARCHAR(250) DEFAULT "",
		       		about TEXT DEFAULT "",
		       		reset_token VARCHAR(250) DEFAULT "",
		       		is_banned INTEGER DEFAULT 0,
				is_superadmin INTEGER DEFAULT 0,
		       		created_date INTEGER,
		       		updated_date INTEGER,
		       		reset_token_date INTEGER DEFAULT 0
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE UNIQUE INDEX users_username_index on users(username);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX users_email_index on users(email);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX users_reset_token_index on users(reset_token);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE groups(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		name VARCHAR(200),
		       		desc TEXT DEFAULT "",
		       		header_msg TEXT DEFAULT "",
		       		is_sticky INTEGER DEFAULT 0,
		       		is_closed INTEGER DEFAULT 0,
		       		created_date INTEGER,
		       		updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX groups_sticky_index on groups(is_sticky);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX groups_name_index on groups(name);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE topics(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title VARCHAR(200) DEFAULT "",
				content TEXT DEFAULT "",
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
				is_deleted INTEGER DEFAULT 0,
				is_sticky INTEGER DEFAULT 0,
				is_closed INTEGER DEFAULT 0,
				num_comments INTEGER DEFAULT 0,
				created_date INTEGER,
				updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topics_userid_index on topics(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topics_groupid_index on topics(groupid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE comments(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT DEFAULT "",
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				topicid INTEGER REFERENCES topics(id) ON DELETE CASCADE,
				parentid INTEGER REFERENCES comments(id) ON DELETE CASCADE,
				is_deleted INTEGER DEFAULT 0,
				is_sticky INTEGER DEFAULT 0,
				created_date INTEGER,
				updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX comments_userid_index on comments(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX comments_topicid_index on comments(topicid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX comments_parentid_index on comments(parentid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE mods(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
		       		created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX mods_userid_index on mods(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX mods_groupid_index on mods(groupid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE admins(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
		       		created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX admins_userid_index on admins(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX admins_groupid_index on admins(groupid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE topicsubscriptions(
				id INTEGER PRIMARY KEY,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				topicid INTEGER REFERENCES topics(id) ON DELETE CASCADE,
				token VARCHAR(128),
				created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topicsubscriptions_userid_index on topicsubscriptions(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topicsubscriptions_topicid_index on topicsubscriptions(topicid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topicsubscriptions_token_index on topicsubscriptions(token);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE groupsubscriptions(
				id INTEGER PRIMARY KEY,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
				token VARCHAR(128),
				created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX groupsubscriptions_userid_index on groupsubscriptions(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX groupsubscriptions_groupid_index on groupsubscriptions(groupid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX groupsubscriptions_token_index on groupsubscriptions(token);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE extranotes(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name VARCHAR(250) NOT NULL,
				content TEXT DEFAULT "",
				URL VARCHAR(250) DEFAULT "",
				created_date INTEGER,
				updated_date INTEGER
	);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE sessions(
				id INTEGER PRIMARY KEY,
				sessionid VARCHAR(250) NOT NULL,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				csrf VARCHAR(250) NOT NULL,
				msg VARCHAR(250) NOT NULL,
				created_date INTEGER NOT NULL,
				updated_date INTEGER NOT NULL
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX sessions_sessionid_index on sessions(sessionid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX sessions_userid_index on sessions(userid);`); err != nil { panic(err) }
}


func DBVersion() int {
	row := db.QueryRow(`SELECT val FROM configs WHERE key="version";`)
	sval := "0"
	if err := row.Scan(&sval); err == nil {
		if ival, err := strconv.Atoi(sval); err == nil {
			return ival
		}
	}
	return 0
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