// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"time"
	"github.com/s-gv/orangeforum/models/db"
	"math/rand"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"encoding/hex"
	"database/sql"
	"net/http"
	"html/template"
	"net/url"
)

const (
	VoteUp = 1
	VoteDown = 2
	VoteFlag = 3
)

const ModelVersion = 1

const (
	ForumName string = "forum_name"
	HeaderMsg string = "header_msg"
	SignupDisabled string = "signup_disabled"
	GroupCreationDisabled string = "group_creation_disabled"
	ImageUploadEnabled string = "image_upload_enabled"
	AllowGroupSubscription string = "allow_group_subscription"
	AllowTopicSubscription string = "allow_topic_subscription"
	DataDir string = "data_dir"
	BodyAppendage string = "body_appendage"
	DefaultFromMail string = "default_from_mail"
	SMTPHost string = "smtp_host"
	SMTPPort string = "smtp_port"
	SMTPUser string = "smtp_user"
	SMTPPass string = "smtp_pass"
)

var ErrIncorrectPasswd = errors.New("Incorrect username/password.")
var ErrUserNotFound = errors.New("Username not found.")
var ErrUserAlreadyExists = errors.New("Username already exists.")

type User struct {
	ID int
	Username string
	PasswdHash string
	Email string
	About string
	IsBanned bool
	IsSuperAdmin bool
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Group struct {
	ID int
	Name string
	Desc string
	IsSticky string
	IsPrivate string
	IsClosed string
	CreatedDate time.Time
	UpdatedDate time.Time
}

type Mod struct {
	ID int
	UserID int
	GroupID int
	CreatedDate time.Time
}

type Admin struct {
	ID int
	UserID int
	GroupID int
	CreatedDate time.Time
}

type CommentVote struct {
	ID int
	UserID int
	CommentID int
	VoteType int
	CreatedDate time.Time
}

type ExtraNote struct {
	ID int
	Name string
	Content string
	URL string
	CreatedDate time.Time
	UpdatedDate time.Time
}

type CommonData struct {
	CSRF string
	Msg string
	UserName string
	IsSuperAdmin bool
	ForumName string
	CurrentURL template.URL
	BodyAppendage string
	IsGroupSubAllowed bool
	IsTopicSubAllowed bool
	ExtraNotesShort []ExtraNote
}

func createUser(userName string, passwd string, email string, isSuperAdmin bool) error {
	if passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err == nil {
		r := db.QueryRow(`SELECT username FROM users WHERE username=?;`, userName)
		var tmp string
		if err := r.Scan(&tmp); err == sql.ErrNoRows {
			db.Exec(`INSERT INTO users(username, passwdhash, email, is_superadmin, created_date, updated_date) VALUES(?, ?, ?, ?, ?, ?);`,
				userName, hex.EncodeToString(passwdHash), email, isSuperAdmin, time.Now().Unix(), time.Now().Unix())
		} else {
			return ErrUserAlreadyExists
		}
	} else {
		return err
	}
	return nil
}

func CreateUser(userName string, passwd string, email string) error {
	return createUser(userName, passwd, email, false)
}

func CreateSuperUser(userName string, passwd string) error {
	return createUser(userName, passwd, "", true)
}

func ReadUserEmail(userName string) string {
	r := db.QueryRow(`SELECT email FROM users WHERE username=?;`, userName)
	var email string
	if err := r.Scan(&email); err == nil {
		return email
	}
	return ""
}

func ReadUserNameByToken(resetToken string) (string, error) {
	if len(resetToken) > 0 {
		r := db.QueryRow(`SELECT username, reset_token_date FROM users WHERE reset_token=?;`, resetToken)
		var userName string
		var rDate int64
		if err := r.Scan(&userName, &rDate); err == nil {
			resetDate := time.Unix(rDate, 0)
			if resetDate.After(time.Now().Add(-48*time.Hour)) {
				return userName, nil
			}
		}
	}
	return "", errors.New("Invalid/Expired reset token.")
}

func ReadUserIDByName(userName string) (int, error) {
	r := db.QueryRow(`SELECT id FROM users WHERE username=?;`, userName)
	var id int
	if err := r.Scan(&id); err == nil {
		return id, nil
	}
	return 0, errors.New("User not found.")
}

func UpdateUserPasswd(userName string, passwd string) error {
	if passwdHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost); err == nil {
		db.Exec(`UPDATE users SET passwdhash=?, reset_token='', reset_token_date=0 WHERE username=?`, hex.EncodeToString(passwdHash), userName)
	} else {
		return err
	}
	return nil
}

func CreateResetToken(userName string) string {
	resetToken := RandSeq(64)
	db.Exec(`UPDATE users SET reset_token=?, reset_token_date=? WHERE username=?;`, resetToken, int64(time.Now().Unix()), userName)
	return resetToken
}

func ProbeUser(userName string) bool {
	r := db.QueryRow(`SELECT username FROM users WHERE username=?;`, userName)
	var tmp string
	if err := r.Scan(&tmp); err == sql.ErrNoRows {
		return false
	}
	return true
}

func ReadGroupIDByName(name string) string {
	r := db.QueryRow(`SELECT id FROM groups WHERE name=?;`, name)
	var id string
	if err := r.Scan(&id); err == nil {
		return id
	}
	return ""
}

func DeleteGroup(groupID string) {
	db.Exec(`UPDATE groups SET is_closed=1 WHERE id=?;`, groupID)
}

func UndeleteGroup(groupID string) {
	db.Exec(`UPDATE groups SET is_closed=0 WHERE id=?;`, groupID)
}

func CreateMod(userName string, groupID string) {
	if uid, err := ReadUserIDByName(userName); err == nil {
		db.Exec(`INSERT INTO mods(userid, groupid, created_date) VALUES(?, ?, ?);`, uid, groupID, time.Now().Unix())
	}
}

func ReadMods(groupID string) []string {
	rows := db.Query(`SELECT users.username FROM users INNER JOIN mods ON users.id=mods.userid WHERE mods.groupid=?;`, groupID)
	var mods []string
	for rows.Next() {
		var mod string
		rows.Scan(&mod)
		mods = append(mods, mod)
	}
	return mods
}

func DeleteMods(groupID string) {
	db.Exec(`DELETE FROM admins WHERE groupid=?;`, groupID)
}


func CreateAdmin(userName string, groupID string) {
	if uid, err := ReadUserIDByName(userName); err == nil {
		db.Exec(`INSERT INTO admins(userid, groupid, created_date) VALUES(?, ?, ?);`, uid, groupID, time.Now().Unix())
	}
}

func ReadAdmins(groupID string) []string {
	rows := db.Query(`SELECT users.username FROM users INNER JOIN admins ON users.id=admins.userid WHERE admins.groupid=?;`, groupID)
	var admins []string
	for rows.Next() {
		var admin string
		rows.Scan(&admin)
		admins = append(admins, admin)
	}
	return admins
}

func IsUserGroupAdmin(userID string, groupID string) bool {
	r := db.QueryRow(`SELECT id FROM admins WHERE userid=? AND groupid=?`, userID, groupID)
	var tmp string
	if err := r.Scan(&tmp); err == nil {
		return true
	}
	return false
}

func DeleteAdmins(groupID string) {
	db.Exec(`DELETE FROM mods WHERE groupid=?;`, groupID)
}



func CreateExtraNote(name string, URL string, content string) {
	db.Exec(`INSERT INTO extranotes(name, URL, content, created_date, updated_date) VALUES(?, ?, ?, ?, ?);`, name, URL, content, time.Now().Unix(), time.Now().Unix())
}

func ReadExtraNotes() []ExtraNote {
	rows := db.Query(`SELECT id, name, URL, content FROM extranotes;`)
	var extraNotes []ExtraNote
	for rows.Next() {
		var extraNote ExtraNote
		rows.Scan(&extraNote.ID, &extraNote.Name, &extraNote.URL, &extraNote.Content)
		extraNotes = append(extraNotes, extraNote)
	}
	return extraNotes
}

func ReadExtraNote(id string) (ExtraNote, error) {
	r := db.QueryRow(`SELECT name, URL, content, created_date, updated_date FROM extranotes WHERE id=?;`, id)
	var e ExtraNote
	var cDate int64
	var uDate int64
	if err := r.Scan(&e.Name, &e.URL, &e.Content, &cDate, &uDate); err == nil {
		e.CreatedDate = time.Unix(cDate, 0)
		e.UpdatedDate = time.Unix(uDate, 0)
		return e, nil
	}
	return ExtraNote{}, errors.New("No note with that ID found")
}

func ReadExtraNotesShort() []ExtraNote {
	rows := db.Query(`SELECT id, name FROM extranotes;`)
	var extraNotes []ExtraNote
	for rows.Next() {
		var extraNote ExtraNote
		rows.Scan(&extraNote.ID, &extraNote.Name)
		extraNotes = append(extraNotes, extraNote)
	}
	return extraNotes
}

func UpdateExtraNote(id string, name string, URL string, content string) {
	now := time.Now()
	db.Exec(`UPDATE extranotes SET name=?, URL=?, content=?, updated_date=? WHERE id=?;`, name, URL, content, int64(now.Unix()), id)
}

func DeleteExtraNote(id string) {
	db.Exec(`DELETE FROM extranotes WHERE id=?;`, id)
}

func ReadCommonData(r *http.Request, sess Session) CommonData {
	userName := ""
	isSuperAdmin := false
	if sess.UserID.Valid {
		r := db.QueryRow(`SELECT username, is_superadmin FROM users WHERE id=?;`, sess.UserID)
		r.Scan(&userName, &isSuperAdmin)
	}
	currentURL := "/"
	if r.URL.Path != "" {
		currentURL = r.URL.Path
		if r.URL.RawQuery != "" {
			currentURL = currentURL + "?" + r.URL.RawQuery
		}
	}
	return CommonData{
		CSRF:sess.CSRFToken,
		Msg:sess.FlashMsg(),
		UserName:userName,
		IsSuperAdmin:isSuperAdmin,
		ForumName:Config(ForumName),
		CurrentURL:template.URL(url.QueryEscape(currentURL)),
		IsGroupSubAllowed:Config(AllowGroupSubscription) != "0",
		IsTopicSubAllowed:Config(AllowTopicSubscription) != "0",
		BodyAppendage:Config(BodyAppendage),
		ExtraNotesShort:ReadExtraNotesShort(),
	}
}

func RandSeq(n int) string {
	var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func WriteConfig(key string, val string) {
	var oldVal string
	if db.QueryRow(`SELECT val FROM configs WHERE name=?;`, key).Scan(&oldVal) == nil {
		if oldVal != val {
			db.Exec(`UPDATE configs SET val=? WHERE name=?;`, val, key)
		}
	} else {
		db.Exec(`INSERT INTO configs(name, val) values(?, ?);`, key, val)
	}
}


func Config(key string) string {
	row := db.QueryRow(`SELECT val FROM configs WHERE name=?;`, key)
	var val string
	if err := row.Scan(&val); err == nil {
		return val
	}
	return "0"
}

func ConfigAllVals() map[string]interface{} {
	vals := map[string]interface{}{
		"forum_name": Config(ForumName),
		"header_msg": Config(HeaderMsg),
		"signup_disabled": Config(SignupDisabled) == "1",
		"group_creation_disabled": Config(GroupCreationDisabled) == "1",
		"image_upload_enabled": Config(ImageUploadEnabled) == "1",
		"allow_group_subscription": Config(AllowGroupSubscription) == "1",
		"allow_topic_subscription": Config(AllowTopicSubscription) == "1",
		"data_dir": Config(DataDir),
		"body_appendage": Config(BodyAppendage),
		"default_from_mail": Config(DefaultFromMail),
		"smtp_host": Config(SMTPHost),
		"smtp_port": Config(SMTPPort),
		"smtp_user": Config(SMTPUser),
		"smtp_pass": Config(SMTPPass),
	}
	return vals
}

func NumUsers() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM users LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func NumGroups() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM groups LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func NumTopics() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM topics LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func NumComments() int64 {
	r := db.QueryRow(`SELECT MAX(_ROWID_) FROM comments LIMIT 1;`)
	var n sql.NullInt64
	if err := r.Scan(&n); err == nil {
		return n.Int64
	}
	return 0
}

func CreateTables() {
	db.Exec(`CREATE TABLE configs(name VARCHAR(250), val TEXT);`)
	db.Exec(`CREATE UNIQUE INDEX configs_key_index on configs(name);`)

	db.Exec(`CREATE TABLE users(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		username VARCHAR(32) NOT NULL,
		       		passwdhash VARCHAR(250) NOT NULL,
		       		email VARCHAR(250) DEFAULT '',
		       		about TEXT DEFAULT '',
		       		reset_token VARCHAR(250) DEFAULT '',
		       		is_banned INTEGER DEFAULT 0,
				is_superadmin INTEGER DEFAULT 0,
		       		created_date INTEGER,
		       		updated_date INTEGER,
		       		reset_token_date INTEGER DEFAULT 0
	);`)
	db.Exec(`CREATE UNIQUE INDEX users_username_index on users(username);`)
	db.Exec(`CREATE INDEX users_email_index on users(email);`)
	db.Exec(`CREATE INDEX users_reset_token_index on users(reset_token);`)
	db.Exec(`CREATE INDEX users_created_index on users(created_date);`)

	db.Exec(`CREATE TABLE groups(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		name VARCHAR(200),
		       		description TEXT DEFAULT '',
		       		header_msg TEXT DEFAULT '',
		       		is_sticky INTEGER DEFAULT 0,
		       		is_closed INTEGER DEFAULT 0,
		       		created_date INTEGER,
		       		updated_date INTEGER
	);`)
	db.Exec(`CREATE INDEX groups_sticky_index on groups(is_sticky);`)
	db.Exec(`CREATE INDEX groups_closed_sticky_index on groups(is_closed, is_sticky DESC);`)
	db.Exec(`CREATE UNIQUE INDEX groups_name_index on groups(name);`)
	db.Exec(`CREATE INDEX groups_created_index on groups(created_date);`)

	db.Exec(`CREATE TABLE topics(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title VARCHAR(200) DEFAULT '',
				content TEXT DEFAULT '',
				image TEXT DEFAULT '',
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
				is_deleted INTEGER DEFAULT 0,
				is_sticky INTEGER DEFAULT 0,
				is_closed INTEGER DEFAULT 0,
				num_comments INTEGER DEFAULT 0,
				created_date INTEGER,
				updated_date INTEGER
	);`)
	db.Exec(`CREATE INDEX topics_userid_created_index on topics(userid, created_date);`)
	db.Exec(`CREATE INDEX topics_groupid_sticky_created_index on topics(groupid, is_sticky DESC, created_date DESC);`)
	db.Exec(`CREATE INDEX topics_created_index on topics(created_date);`)

	db.Exec(`CREATE TABLE comments(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				content TEXT DEFAULT '',
				image TEXT DEFAULT '',
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				topicid INTEGER REFERENCES topics(id) ON DELETE CASCADE,
				parentid INTEGER REFERENCES comments(id) ON DELETE CASCADE,
				is_deleted INTEGER DEFAULT 0,
				is_sticky INTEGER DEFAULT 0,
				created_date INTEGER,
				updated_date INTEGER
	);`)
	db.Exec(`CREATE INDEX comments_userid_created_index on comments(userid, created_date);`)
	db.Exec(`CREATE INDEX comments_parentid_index on comments(parentid);`)
	db.Exec(`CREATE INDEX comments_topicid_sticky_created_index on comments(topicid, is_sticky DESC, created_date);`)
	db.Exec(`CREATE INDEX comments_created_index on comments(created_date);`)

	db.Exec(`CREATE TABLE mods(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
		       		created_date INTEGER
	);`)
	db.Exec(`CREATE INDEX mods_userid_index on mods(userid);`)
	db.Exec(`CREATE INDEX mods_groupid_userid_index on mods(groupid, userid);`)

	db.Exec(`CREATE TABLE admins(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
		       		userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
		       		created_date INTEGER
	);`)
	db.Exec(`CREATE INDEX admins_userid_index on admins(userid);`)
	db.Exec(`CREATE INDEX admins_groupid_userid_index on admins(groupid, userid);`)

	db.Exec(`CREATE TABLE topicsubscriptions(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				topicid INTEGER REFERENCES topics(id) ON DELETE CASCADE,
				token VARCHAR(128),
				created_date INTEGER
	);`)
	db.Exec(`CREATE INDEX topicsubscriptions_userid_index on topicsubscriptions(userid);`)
	db.Exec(`CREATE INDEX topicsubscriptions_topicid_userid_index on topicsubscriptions(topicid, userid);`)
	db.Exec(`CREATE INDEX topicsubscriptions_token_index on topicsubscriptions(token);`)

	db.Exec(`CREATE TABLE groupsubscriptions(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				groupid INTEGER REFERENCES groups(id) ON DELETE CASCADE,
				token VARCHAR(128),
				created_date INTEGER
	);`)
	db.Exec(`CREATE INDEX groupsubscriptions_userid_index on groupsubscriptions(userid);`)
	db.Exec(`CREATE INDEX groupsubscriptions_groupid_userid_index on groupsubscriptions(groupid, userid);`)
	db.Exec(`CREATE INDEX groupsubscriptions_token_index on groupsubscriptions(token);`)

	db.Exec(`CREATE TABLE extranotes(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name VARCHAR(250) NOT NULL,
				content TEXT DEFAULT '',
				URL VARCHAR(250) DEFAULT '',
				created_date INTEGER,
				updated_date INTEGER
	);`)

	db.Exec(`CREATE TABLE sessions(
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				sessionid VARCHAR(250) NOT NULL,
				userid INTEGER REFERENCES users(id) ON DELETE CASCADE,
				csrf VARCHAR(250) NOT NULL,
				msg VARCHAR(250) NOT NULL,
				created_date INTEGER NOT NULL,
				updated_date INTEGER NOT NULL
	);`)
	db.Exec(`CREATE INDEX sessions_sessionid_index on sessions(sessionid);`)
	db.Exec(`CREATE INDEX sessions_userid_index on sessions(userid);`)
}

func Migrate() {
	CreateTables()

	WriteConfig("version", "1");
	WriteConfig(HeaderMsg, "")
	WriteConfig(ForumName, "Orange Forum")
	WriteConfig(SignupDisabled, "0")
	WriteConfig(GroupCreationDisabled, "0")
	WriteConfig(ImageUploadEnabled, "0")
	WriteConfig(AllowGroupSubscription, "0")
	WriteConfig(AllowTopicSubscription, "0")
	WriteConfig(DataDir, "")
	WriteConfig(BodyAppendage, "")
	WriteConfig(DefaultFromMail, "admin@example.com")
	WriteConfig(SMTPHost, "")
	WriteConfig(SMTPPort, "25")
	WriteConfig(SMTPUser, "")
	WriteConfig(SMTPPass, "")
}

func IsMigrationNeeded() bool {
	dbver := db.Version()
	return dbver != ModelVersion
}
