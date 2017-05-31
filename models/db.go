package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"log"
	"time"
	"fmt"
	"math/rand"
)

type ConfigKey string

const (
	Version ConfigKey = "version"
	Secret ConfigKey = "secret"
	ForumName ConfigKey = "title"
	HeaderMsg ConfigKey = "header_msg"
	SignupNeedsApproval ConfigKey = "signup_needs_approval"
	PublicViewDisabled ConfigKey = "public_view_disabled"
	SignupDisabled ConfigKey = "signup_disabled"
	ImageUploadEnabled ConfigKey = "image_upload_enabled"
	FileUploadEnabled ConfigKey = "file_upload_enabled"
	AllowGroupSubscription ConfigKey = "allow_group_subscription"
	AllowTopicSubscription ConfigKey = "allow_topic_subscription"
	AutoSubscribeToMyTopic ConfigKey = "auto_subscribe_to_my_topic"
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

	if _, err := db.Exec(`CREATE TABLE users(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		username TEXT NOT NULL,
		       		passwdhash TEXT,
		       		email TEXT,
		       		about TEXT,
		       		karma INTEGER,
		       		is_banned BOOLEAN,
		       		is_warned BOOLEAN,
				is_superadmin BOOLEAN,
				is_supermod BOOLEAN,
				is_approved BOOLEAN,
				secret TEXT,
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
	if _, err := db.Exec(`CREATE UNIQUE INDEX groups_sticky_index on groups(is_sticky);`); err != nil { panic(err) }
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


	WriteConfig(Version, "1")
	WriteConfig(HeaderMsg, "")
	WriteConfig(ForumName, "OrangeForum")
	WriteConfig(Secret, randSeq(32))
	WriteConfig(SignupNeedsApproval, "0")
	WriteConfig(PublicViewDisabled, "0")
	WriteConfig(SignupDisabled, "0")
	WriteConfig(FileUploadEnabled, "0")
	WriteConfig(ImageUploadEnabled, "0")
	WriteConfig(AllowGroupSubscription, "0")
	WriteConfig(AllowTopicSubscription, "0")
	WriteConfig(AutoSubscribeToMyTopic, "0")
}



func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


func DBVersion() int {
	if val, err := strconv.Atoi(Config("version", "0")); err == nil {
		return val
	}
	return 0

}

func WriteConfig(key ConfigKey, val string) error {
	return exec("WriteConfig", `INSERT OR REPLACE INTO config(key, val) values(?, ?);`, key, val)
}


func Config(key ConfigKey, defaultVal string) string {
	row, err := queryRow("Config", `SELECT val FROM config WHERE key=?;`, key)
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

	rand.Seed(time.Now().UnixNano())

	dbver := DBVersion()
	if dbver < ModelVersion {
		return ErrDBVer
	}
	return nil
}

func Benchmark() {
	start := time.Now()
	for i := 0; i < 100000; i++ {
		x := Config("version", "0")//WriteConfig("version", "3")
		if x == "0" {
			panic("Er")
		}
	}
	elapsed := time.Since(start)
	WriteConfig("version", "2")
	println("Test val:", Config("version", "0"))
	fmt.Printf("time: %s\n", elapsed)
}