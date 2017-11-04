// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import (
	"github.com/s-gv/orangeforum/models/db"
	"log"
)

const ModelVersion = 2

func Migration1() {
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
	// db.Exec(`ALTER TABLE groups ADD COLUMN is_private INTEGER DEFAULT 0;`) // Migration 2.
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
	// db.Exec(`ALTER TABLE topics ADD COLUMN activity_date INTEGER;`) // Migration 2. Default value set to created_date
	db.Exec(`CREATE INDEX topics_userid_created_index on topics(userid, created_date);`)
	db.Exec(`CREATE INDEX topics_groupid_sticky_created_index on topics(groupid, is_sticky DESC, created_date DESC);`)
	db.Exec(`CREATE INDEX topics_created_index on topics(created_date);`)
	// db.Exec(`CREATE INDEX topics_groupid_sticky_activity_index on topics(groupid, is_sticky DESC, activity_date DESC);`) // Migration 2
	// db.Exec(`CREATE INDEX topics_activity_index on topics(activity_date);`) // Migration 2

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

func Migration2() {
	db.Exec(`ALTER TABLE topics ADD COLUMN activity_date INTEGER;`)
	db.Exec(`UPDATE topics SET activity_date = created_date;`)
	db.Exec(`CREATE INDEX topics_groupid_sticky_activity_index on topics(groupid, is_sticky DESC, activity_date DESC);`)
	db.Exec(`CREATE INDEX topics_activity_index on topics(activity_date);`)

	db.Exec(`ALTER TABLE groups ADD COLUMN is_private INTEGER DEFAULT 0;`)
}

func Migrate() {
	dbver := db.Version()
	if dbver == ModelVersion {
		log.Panicf("[ERROR] DB migration not needed. DB up-to-date.\n")
	} else if dbver > ModelVersion {
		log.Panicf("[ERROR] DB version (%d) is greater than binary version (%d). Use newer binary.\n", dbver, ModelVersion)
	}
	for dbver < ModelVersion {
		if dbver == 0 {
			Migration1()

			WriteConfig(Version, "1");
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
		} else if dbver == 1 {
			Migration2()

			WriteConfig(Version, "2");
		}
		dbver = db.Version()
	}
}
