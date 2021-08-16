// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import "github.com/jmoiron/sqlx"

func migrate001(db *sqlx.DB) {
	db.MustExec(`CREATE OR REPLACE FUNCTION update_modified_timestamp() RETURNS TRIGGER LANGUAGE plpgsql AS
		$$
		BEGIN
			new.updated_at := current_timestamp;
			RETURN new;
		END;
		$$;
	`)

	db.MustExec(`CREATE TABLE configs(
		name    VARCHAR(250) NOT NULL PRIMARY KEY,
		val     VARCHAR(250) NOT NULL DEFAULT ''
	);`)

	db.MustExec(`CREATE TABLE domains(
		domain_id                   SERIAL PRIMARY KEY,
		domain_name                 VARCHAR(250) NOT NULL UNIQUE,
		forum_name                  VARCHAR(250) NOT NULL DEFAULT 'Orange Forum',
		no_regular_signup_msg       VARCHAR(250) NOT NULL DEFAULT '',
		signup_token                VARCHAR(250) NOT NULL DEFAULT '',
		edit_window                 INTEGER DEFAULT 20,
		auto_topic_close_days       INTEGER DEFAULT 60,
		user_activity_window        INTEGER DEFAULT 3,
		max_num_activity            INTEGER DEFAULT 20,
		is_regular_signup_enabled   BOOL NOT NULL DEFAULT false,
		is_readonly                 BOOL NOT NULL DEFAULT false,
		enable_group_sub            BOOL NOT NULL DEFAULT false,
		enable_topic_autosub        BOOL NOT NULL DEFAULT false,
		enable_comment_autosub      BOOL NOT NULL DEFAULT false,
		archived_at                 TIMESTAMPTZ,
		created_at                  TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		updated_at                  TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER update_timestamp BEFORE UPDATE ON domains FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();`)
	db.MustExec(`CREATE UNIQUE INDEX domains_domain_index ON domains(domain_name);`)

	db.MustExec(`CREATE TABLE users(
		user_id                             SERIAL PRIMARY KEY,
		domain_id                           INTEGER NOT NULL REFERENCES domains(domain_id) ON DELETE CASCADE,
		email                               VARCHAR(250) NOT NULL,
		display_name                        VARCHAR(32) NOT NULL,
		passwd_hash                         VARCHAR(250) NOT NULL,
		about                               TEXT NOT NULL DEFAULT '',
		is_superadmin                       BOOL NOT NULL DEFAULT false,
		is_supermod                         BOOL NOT NULL DEFAULT false,
		is_topic_autosubscribe              BOOL NOT NULL DEFAULT true,
		is_comment_autosubscribe            BOOL NOT NULL DEFAULT true,
		is_email_notifications_disabled     BOOL NOT NULL DEFAULT false,
		num_topics                          INTEGER NOT NULL DEFAULT 0,
		num_comments                        INTEGER NOT NULL DEFAULT 0,
		num_activity                        INTEGER NOT NULL DEFAULT 0,
		onetime_login_token                 VARCHAR(250) NOT NULL DEFAULT '',
		onetime_login_token_at              TIMESTAMPTZ NOT NULL DEFAULT to_timestamp(0),
		reset_token                         VARCHAR(250) NOT NULL DEFAULT '',
		last_ip                             VARCHAR(50) NOT NULL DEFAULT '',
		activity_at                         TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		reset_at                            TIMESTAMPTZ NOT NULL DEFAULT to_timestamp(0),
		logout_at                           TIMESTAMPTZ NOT NULL DEFAULT to_timestamp(0),
		banned_at                           TIMESTAMPTZ,
		archived_at                         TIMESTAMPTZ,
		created_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		updated_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER update_timestamp BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();`)
	db.MustExec(`CREATE UNIQUE INDEX users_domain_email_index  ON users(domain_id, email);`)
	db.MustExec(`CREATE INDEX users_supermod_index             ON users(domain_id, is_supermod);`)
	db.MustExec(`CREATE INDEX users_otp_token_index            ON users(onetime_login_token);`)
	db.MustExec(`CREATE INDEX users_reset_token_index          ON users(reset_token);`)
	db.MustExec(`CREATE INDEX users_created_index              ON users(created_at);`)

	db.MustExec(`CREATE TABLE categories(
		category_id                         SERIAL PRIMARY KEY,
		domain_id                           INTEGER NOT NULL REFERENCES domains(domain_id) ON DELETE CASCADE,
		name                                VARCHAR(250) NOT NULL,
		description                         VARCHAR(250) NOT NULL,
		is_private                          BOOL NOT NULL DEFAULT false,
		is_readonly                         BOOL NOT NULL DEFAULT false,
		is_restricted                       BOOL NOT NULL DEFAULT false,
		archived_at                         TIMESTAMPTZ,
		created_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		updated_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER update_timestamp BEFORE UPDATE ON categories FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();`)
	db.MustExec(`CREATE INDEX categories_domain_index          ON categories(domain_id);`)

	db.MustExec(`CREATE TABLE topics(
		topic_id                            SERIAL PRIMARY KEY,
		category_id                         INTEGER NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
		user_id                             INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
		title                               VARCHAR(250) NOT NULL,
		content                             TEXT NOT NULL DEFAULT '',
		is_sticky                           BOOL NOT NULL DEFAULT false,
		is_readonly                         BOOL NOT NULL DEFAULT false,
		num_comments                        INTEGER NOT NULL DEFAULT 0,
		num_views                           INTEGER NOT NULL DEFAULT 0,
		activity_at                         TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		archived_at                         TIMESTAMPTZ,
		created_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		updated_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER update_timestamp BEFORE UPDATE          ON topics FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();`)
	db.MustExec(`CREATE INDEX topics_category_sticky_activity_index     ON topics(category_id, is_sticky, activity_at);`)

	db.MustExec(`CREATE TABLE comments(
		comment_id                          SERIAL PRIMARY KEY,
		topic_id                            INTEGER NOT NULL REFERENCES topics(topic_id) ON DELETE CASCADE,
		user_id                             INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
		content                             TEXT NOT NULL DEFAULT '',
		is_sticky                           BOOL NOT NULL DEFAULT false,
		archived_at                         TIMESTAMPTZ,
		created_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		updated_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER update_timestamp BEFORE UPDATE          ON comments FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();`)
	db.MustExec(`CREATE INDEX comments_topic_sticky_created_index       ON comments(topic_id, is_sticky, created_at);`)

	db.MustExec(`CREATE TABLE category_subs(
		sub_id                              SERIAL PRIMARY KEY,
		category_id                         INTEGER NOT NULL REFERENCES categories(category_id) ON DELETE CASCADE,
		user_id                             INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
		unsub_token                         VARCHAR(64) NOT NULL DEFAULT '',
		created_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE INDEX catsubs_cat_index               ON category_subs(category_id);`)
	db.MustExec(`CREATE INDEX catsubs_unsub_token_index       ON category_subs(unsub_token);`)

	db.MustExec(`CREATE TABLE topic_subs(
		sub_id                              SERIAL PRIMARY KEY,
		topic_id                            INTEGER NOT NULL REFERENCES topics(topic_id) ON DELETE CASCADE,
		user_id                             INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
		unsub_token                         VARCHAR(64) NOT NULL DEFAULT '',
		created_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE INDEX topicsubs_cat_index               ON topic_subs(topic_id);`)
	db.MustExec(`CREATE INDEX topicsubs_unsub_token_index       ON topic_subs(unsub_token);`)

	db.MustExec(`CREATE TABLE notes(
		note_id                             SERIAL PRIMARY KEY,
		domain_id                           INTEGER NOT NULL REFERENCES domains(domain_id) ON DELETE CASCADE,
		name                                VARCHAR(64) NOT NULL DEFAULT '',
		content                             TEXT NOT NULL DEFAULT '',
		url                                 VARCHAR(64) NOT NULL DEFAULT '',
		created_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		updated_at                          TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER update_timestamp BEFORE UPDATE          ON notes FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();`)
	db.MustExec(`CREATE INDEX notes_domain_url_index                    ON notes(domain_id, url);`)

	// Add some config data
	db.MustExec(`INSERT INTO configs(name, val) VALUES('` + DBVersion + `', '1');`)
}
