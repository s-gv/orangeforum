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
		name 	VARCHAR(250) NOT NULL PRIMARY KEY,
		val 	VARCHAR(250) NOT NULL DEFAULT ''
	);`)

	db.MustExec(`CREATE TABLE domains(
		domain_id 					SERIAL PRIMARY KEY,
		domain_name 				VARCHAR(250) NOT NULL UNIQUE,
		forum_name 					VARCHAR(250) NOT NULL DEFAULT '',
		no_regular_signup_msg 		VARCHAR(250) NOT NULL DEFAULT '',
		signup_token 				VARCHAR(250) NOT NULL DEFAULT '',
		edit_window 				INTEGER DEFAULT 20,
		auto_topic_close_days 		INTEGER DEFAULT 60,
		user_activity_window 		INTEGER DEFAULT 3,
		max_num_activity 			INTEGER DEFAULT 20,
		is_default 					BOOL NOT NULL DEFAULT false,
		is_regular_signup_enabled 	BOOL NOT NULL DEFAULT false,
		is_readonly 				BOOL NOT NULL DEFAULT false,
		enable_group_sub 			BOOL NOT NULL DEFAULT false,
		enable_topic_autosub 		BOOL NOT NULL DEFAULT false,
		enable_comment_autosub 		BOOL NOT NULL DEFAULT false,
		archived_at 				TIMESTAMPTZ,
		created_at 					TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		updated_at 					TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER update_timestamp BEFORE UPDATE ON domains FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();`)
	db.MustExec(`CREATE UNIQUE INDEX domains_domain_index ON domains(domain_name);`)

	db.MustExec(`CREATE TABLE users(
		user_id 							SERIAL PRIMARY KEY,
		domain_id 							INTEGER REFERENCES domains(domain_id) ON DELETE CASCADE,
		email 								VARCHAR(250) DEFAULT '',
		username 							VARCHAR(32) NOT NULL UNIQUE,
		passwd_hash 						VARCHAR(250) NOT NULL,
		about 								TEXT NOT NULL DEFAULT '',
		is_superadmin 						BOOL NOT NULL DEFAULT false,
		is_topic_autosubscribe				BOOL NOT NULL DEFAULT true,
		is_comment_autosubscribe			BOOL NOT NULL DEFAULT true,
		is_email_notifications_disabled 	BOOL NOT NULL DEFAULT false,
		num_topics							INTEGER NOT NULL DEFAULT 0,
		num_comments						INTEGER NOT NULL DEFAULT 0,
		num_activity						INTEGER NOT NULL DEFAULT 0,
		onetime_login_token 				VARCHAR(250) NOT NULL DEFAULT '',
		onetime_login_token_at 				TIMESTAMPTZ NOT NULL DEFAULT to_timestamp(0),
		reset_token 						VARCHAR(250) NOT NULL DEFAULT '',
		last_ip								VARCHAR(50) NOT NULL DEFAULT '',
		activity_at							TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		reset_at 							TIMESTAMPTZ NOT NULL DEFAULT to_timestamp(0),
		logout_at 							TIMESTAMPTZ NOT NULL DEFAULT to_timestamp(0),
		banned_at 							TIMESTAMPTZ,
		archived_at 						TIMESTAMPTZ,
		created_at 							TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		updated_at 							TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);`)
	db.MustExec(`CREATE TRIGGER update_timestamp BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();`)
	db.MustExec(`CREATE UNIQUE INDEX users_domain_username_index 	ON users(domain_id, username);`)
	db.MustExec(`CREATE INDEX users_domain_email_index		 		ON users(domain_id, email);`)
	db.MustExec(`CREATE INDEX users_otp_token_index 				ON users(onetime_login_token);`)
	db.MustExec(`CREATE INDEX users_reset_token_index 				ON users(reset_token);`)
	db.MustExec(`CREATE INDEX users_created_index 					ON users(created_at);`)

}
