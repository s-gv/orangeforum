package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"log"
)


var db *sql.DB

var stmts = make(map[string]*sql.Stmt)


func CreateTables() {
	if _, err := db.Exec(`CREATE TABLE configs(key TEXT, val TEXT);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE UNIQUE INDEX configs_key_index on configs(key);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE users(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		username TEXT NOT NULL,
		       		passwdhash TEXT NOT NULL,
		       		email TEXT DEAFULT "",
		       		about TEXT DEFAULT "",
		       		karma INTEGER DEFAULT 0,
		       		is_banned BOOLEAN DEFAULT false,
		       		is_warned BOOLEAN DEFAULT false,
				is_superadmin BOOLEAN DEFAULT false,
				is_supermod BOOLEAN DEFAULT false,
				is_approved BOOLEAN DEFAULT false,
		       		created_date INTEGER,
		       		updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE UNIQUE INDEX users_username_index on users(username);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX users_email_index on users(email);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE groups(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		name TEXT,
		       		desc TEXT,
		       		is_sticky BOOLEAN,
		       		is_private BOOLEAN,
		       		is_closed BOOLEAN,
		       		header_msg TEXT,
		       		created_date INTEGER,
		       		updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX groups_sticky_index on groups(is_sticky);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX groups_created_index on groups(created_date);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE topics(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT,
				authorid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
				is_deleted BOOLEAN,
				is_sticky BOOLEAN,
				is_closed BOOLEAN,
				numcomments INTEGER,
				upvotes INTEGER,
				downvotes INTEGER,
				flagvotes INTEGER,
				created_date INTEGER,
				updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topics_authorid_index on topics(authorid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topics_groupid_index on topics(groupid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE comments(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT,
				authorid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				topicid INTEGER REFERENCES topics(id) ON DELETE CASCADE,
				parentid INTEGER REFERENCES comments(id) ON DELETE CASCADE,
				is_deleted BOOLEAN,
				is_sticky BOOLEAN,
				upvotes INTEGER,
				downvotes INTEGER,
				flagvotes INTEGER,
				created_date INTEGER,
				updated_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX comments_authorid_index on comments(authorid);`); err != nil { panic(err) }
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


	if _, err := db.Exec(`CREATE TABLE topicvotes(
				id INTEGER PRIMARY KEY,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				topicid INTEGER REFERENCES topics(id) ON DELETE CASCADE,
				votetype INTEGER,
				created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topicvotes_userid_index on topicvotes(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topicvotes_topicid_index on topicvotes(topicid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE commentvotes(
				id INTEGER PRIMARY KEY,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				commentid INTEGER REFERENCES comments(id) ON DELETE CASCADE,
				votetype INTEGER,
				created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX commentvotes_userid_index on commentvotes(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX commentvotes_commentid_index on commentvotes(commentid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE topicsubscriptions(
				id INTEGER PRIMARY KEY,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				topicid INTEGER REFERENCES topics(id) ON DELETE CASCADE,
				created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topicsubscriptions_userid_index on topicsubscriptions(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX topicsubscriptions_topicid_index on topicsubscriptions(topicid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE groupsubscriptions(
				id INTEGER PRIMARY KEY,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
				created_date INTEGER
	);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX groupsubscriptions_userid_index on groupsubscriptions(userid);`); err != nil { panic(err) }
	if _, err := db.Exec(`CREATE INDEX groupsubscriptions_groupid_index on groupsubscriptions(groupid);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE extranotes(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL,
				content TEXT,
				URL TEXT,
				created_date INTEGER,
				updated_date INTEGER
	);`); err != nil { panic(err) }


	if _, err := db.Exec(`CREATE TABLE sessions(
				id INTEGER PRIMARY KEY,
				sessionid TEXT NOT NULL,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				csrf TEXT NOT NULL,
				msg TEXT NOT NULL,
				data TEXT NOT NULL,
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

func QueryRow(query string, args ...interface{}) *sql.Row {
	stmt := makeStmt(query)
	return stmt.QueryRow(args...)
}

func ScanRow(row *sql.Row, args ...interface{}) error {
	err := row.Scan(args...)
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