// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import "github.com/jmoiron/sqlx"

func migrateSqlite0001(db *sqlx.DB) {
	db.MustExec(`CREATE TABLE configs(name VARCHAR(250) NOT NULL, val TEXT NOT NULL DEFAULT '');`)
	db.MustExec(`CREATE UNIQUE INDEX configs_key_index on configs(name);`)

	db.MustExec(`CREATE TABLE users(
		id UUID NOT NULL PRIMARY KEY,
		email VARCHAR(250) NOT NULL DEFAULT '',
		username VARCHAR(32) NOT NULL,
		passwd_hash VARCHAR(250) NOT NULL,
		about TEXT NOT NULL DEFAULT '',
		is_superadmin BOOL NOT NULL DEFAULT false,
		is_email_notifications_disabled BOOL NOT NULL DEFAULT false,
		onetime_login_token VARCHAR(250) NOT NULL DEFAULT '',
		onetime_login_at DATETIME NOT NULL DEFAULT (datetime(0, 'unixepoch')),
		reset_token VARCHAR(250) NOT NULL DEFAULT '',
		reset_at DATETIME NOT NULL DEFAULT (datetime(0, 'unixepoch')),
		logout_at DATETIME NOT NULL DEFAULT (datetime(0, 'unixepoch')),
		banned_at DATETIME,
		archived_at DATETIME,
		created_at DATETIME NOT NULL DEFAULT current_timestamp,
		updated_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER users_update_timestamp BEFORE UPDATE ON users
		FOR EACH ROW BEGIN
			UPDATE users SET updated_at = current_timestamp WHERE id = NEW.id;
		END;
	`)
	db.MustExec(`CREATE UNIQUE INDEX users_username_index on users(username);`)
	db.MustExec(`CREATE UNIQUE INDEX users_email_index on users(email);`)
	db.MustExec(`CREATE INDEX users_otp_token_index on users(onetime_login_token);`)
	db.MustExec(`CREATE INDEX users_reset_token_index on users(reset_token);`)
	db.MustExec(`CREATE INDEX users_created_index on users(created_at);`)

	db.MustExec(`CREATE TABLE groups(
		id UUID NOT NULL PRIMARY KEY,
		name VARCHAR(200) NOT NULL DEFAULT '',
		description TEXT NOT NULL DEFAULT '',
		header_msg TEXT NOT NULL DEFAULT '',
		is_sticky BOOL NOT NULL DEFAULT false,
		is_restricted BOOL NOT NULL DEFAULT false,
		is_private BOOL NOT NULL DEFAULT false,
		is_readonly BOOL NOT NULL DEFAULT false,
		archived_at DATETIME,
		created_at DATETIME NOT NULL DEFAULT current_timestamp,
		updated_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER groups_update_timestamp BEFORE UPDATE ON groups
		FOR EACH ROW BEGIN
			UPDATE groups SET updated_at = current_timestamp WHERE id = NEW.id;
		END;
	`)
	db.MustExec(`CREATE INDEX groups_sticky_index on groups(is_sticky);`)
	db.MustExec(`CREATE INDEX groups_private_sticky_index on groups(is_private, is_sticky DESC);`)
	db.MustExec(`CREATE UNIQUE INDEX groups_name_index on groups(name);`)
	db.MustExec(`CREATE INDEX groups_created_index on groups(created_at);`)

	db.MustExec(`CREATE TABLE topics(
		id UUID NOT NULL PRIMARY KEY,
		title VARCHAR(200) NOT NULL DEFAULT '',
		content TEXT NOT NULL DEFAULT '',
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		is_sticky BOOL NOT NULL DEFAULT false,
		is_readonly BOOL NOT NULL DEFAULT false,
		num_comments INTEGER NOT NULL DEFAULT 0,
		activity_at DATETIME NOT NULL DEFAULT current_timestamp,
		archived_at DATETIME,
		created_at DATETIME NOT NULL DEFAULT current_timestamp,
		updated_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER topics_update_timestamp BEFORE UPDATE ON topics
		FOR EACH ROW BEGIN
			UPDATE topics SET updated_at = current_timestamp WHERE id = NEW.id;
		END;
	`)
	db.MustExec(`CREATE INDEX topics_userid_created_index on topics(user_id, created_at);`)
	db.MustExec(`CREATE INDEX topics_groupid_sticky_created_index on topics(group_id, is_sticky DESC, created_at DESC);`)
	db.MustExec(`CREATE INDEX topics_created_index on topics(created_at);`)
	db.MustExec(`CREATE INDEX topics_groupid_sticky_activity_index on topics(group_id, is_sticky DESC, activity_at DESC);`)
	db.MustExec(`CREATE INDEX topics_activity_index on topics(activity_at);`)

	db.MustExec(`CREATE TABLE comments(
		id UUID NOT NULL PRIMARY KEY,
		content TEXT NOT NULL DEFAULT '',
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
		is_sticky BOOL NOT NULL DEFAULT false,
		archived_at DATETIME,
		created_at DATETIME NOT NULL DEFAULT current_timestamp,
		updated_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER comments_update_timestamp BEFORE UPDATE ON comments
		FOR EACH ROW BEGIN
			UPDATE comments SET updated_at = current_timestamp WHERE id = NEW.id;
		END;
	`)
	db.MustExec(`CREATE INDEX comments_userid_created_index on comments(user_id, created_at);`)
	db.MustExec(`CREATE INDEX comments_topicid_sticky_created_index on comments(topic_id, is_sticky, created_at);`)
	db.MustExec(`CREATE INDEX comments_topicid_created_index on comments(topic_id, created_at);`)
	db.MustExec(`CREATE INDEX comments_created_index on comments(created_at);`)

	db.MustExec(`CREATE TABLE sub_comments(
		id UUID NOT NULL PRIMARY KEY,
		content TEXT NOT NULL DEFAULT '',
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		comment_id UUID NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
		archived_at DATETIME,
		created_at DATETIME NOT NULL DEFAULT current_timestamp,
		updated_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER sub_comments_update_timestamp BEFORE UPDATE ON sub_comments
		FOR EACH ROW BEGIN
			UPDATE sub_comments SET updated_at = current_timestamp WHERE id = NEW.id;
		END;
	`)
	db.MustExec(`CREATE INDEX sub_comments_userid_created_index on sub_comments(user_id, created_at);`)
	db.MustExec(`CREATE INDEX sub_comments_commentid_created_index on sub_comments(comment_id, created_at);`)
	db.MustExec(`CREATE INDEX sub_comments_created_index on sub_comments(created_at);`)

	db.MustExec(`CREATE TABLE mods(
		id UUID NOT NULL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		created_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE INDEX mods_userid_index on mods(user_id);`)
	db.MustExec(`CREATE INDEX mods_groupid_userid_index on mods(group_id, user_id);`)

	db.MustExec(`CREATE TABLE admins(
		id UUID NOT NULL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		created_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE INDEX admins_userid_index on admins(user_id);`)
	db.MustExec(`CREATE INDEX admins_groupid_userid_index on admins(group_id, user_id);`)

	db.MustExec(`CREATE TABLE whitelisted_users(
		id UUID NOT NULL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		is_whitelisted BOOL NOT NULL DEFAULT false,
		created_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE INDEX whitelisted_users_userid_index on whitelisted_users(user_id);`)
	db.MustExec(`CREATE INDEX whitelisted_users_groupid_userid_index on whitelisted_users(group_id, user_id);`)

	db.MustExec(`CREATE TABLE topic_subs(
		id UUID NOT NULL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		topic_id UUID NOT NULL REFERENCES topics(id) ON DELETE CASCADE,
		unsub_token VARCHAR(128),
		created_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE INDEX topicsubs_userid_index on topic_subs(user_id);`)
	db.MustExec(`CREATE INDEX topicsubs_topicid_userid_index on topic_subs(topic_id, user_id);`)
	db.MustExec(`CREATE INDEX topicsubs_token_index on topic_subs(unsub_token);`)

	db.Exec(`CREATE TABLE group_subs(
		id UUID NOT NULL PRIMARY KEY,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
		unsub_token VARCHAR(128),
		created_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.Exec(`CREATE INDEX groupsubs_userid_index on group_subs(user_id);`)
	db.Exec(`CREATE INDEX groupsubs_groupid_userid_index on group_subs(group_id, userid);`)
	db.Exec(`CREATE INDEX groupsubs_token_index on group_subs(unsub_token);`)

	db.MustExec(`CREATE TABLE extra_notes(
		id UUID NOT NULL PRIMARY KEY,
		name VARCHAR(250) NOT NULL,
		content TEXT NOT NULL DEFAULT '',
		url VARCHAR(250) NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT current_timestamp,
		updated_at DATETIME NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER extranotes_update_timestamp BEFORE UPDATE ON extra_notes
		FOR EACH ROW BEGIN
			UPDATE extranotes SET updated_at = current_timestamp WHERE id = NEW.id;
		END;
	`)
}
