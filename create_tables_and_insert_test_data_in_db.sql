DROP TABLE IF EXISTS configs;
DROP TABLE IF EXISTS bannedips;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS domains;

CREATE OR REPLACE FUNCTION update_modified_timestamp() RETURNS TRIGGER LANGUAGE plpgsql AS
		$$
		BEGIN
			new.updated_at := current_timestamp;
			RETURN new;
		END;
		$$;

CREATE TABLE domains(
		domain_id 					SERIAL PRIMARY KEY,
		domain_name 				VARCHAR(250) NOT NULL UNIQUE,
		forum_name 					VARCHAR(250) NOT NULL DEFAULT 'Orange Forum',
		no_regular_signup_msg 		VARCHAR(250) NOT NULL DEFAULT '',
		signup_token 				VARCHAR(250) NOT NULL DEFAULT '',
		edit_window 				INTEGER DEFAULT 20,
		auto_topic_close_days 		INTEGER DEFAULT 60,
		user_activity_window 		INTEGER DEFAULT 3,
		max_num_activity 			INTEGER DEFAULT 20,
		is_regular_signup_enabled 	BOOL NOT NULL DEFAULT false,
		is_readonly 				BOOL NOT NULL DEFAULT false,
		enable_group_sub 			BOOL NOT NULL DEFAULT false,
		enable_topic_autosub 		BOOL NOT NULL DEFAULT false,
		enable_comment_autosub 		BOOL NOT NULL DEFAULT false,
		archived_at 				TIMESTAMPTZ,
		created_at 					TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
		updated_at 					TIMESTAMPTZ NOT NULL DEFAULT current_timestamp);

CREATE TRIGGER update_timestamp BEFORE UPDATE ON domains FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();
CREATE UNIQUE INDEX domains_domain_index ON domains(domain_name);

CREATE TABLE users(
		user_id 							SERIAL PRIMARY KEY,
		domain_id 							INTEGER NOT NULL REFERENCES domains(domain_id) ON DELETE CASCADE,
		email 								VARCHAR(250) NOT NULL,
		username 							VARCHAR(32) NOT NULL,
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
		updated_at 							TIMESTAMPTZ NOT NULL DEFAULT current_timestamp);
        
CREATE TRIGGER update_timestamp BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_modified_timestamp();
CREATE UNIQUE INDEX users_domain_username_index 	ON users(domain_id, username);
CREATE UNIQUE INDEX users_domain_email_index		ON users(domain_id, email);
CREATE INDEX users_otp_token_index 				ON users(onetime_login_token);
CREATE INDEX users_reset_token_index 				ON users(reset_token);
CREATE INDEX users_created_index 					ON users(created_at);


CREATE TABLE configs(
		name 	VARCHAR(250) NOT NULL PRIMARY KEY,
		val 	VARCHAR(250) NOT NULL DEFAULT '');

INSERT INTO configs(name, val) VALUES('` + DBVersion + `', '1');

CREATE TABLE bannedips (
		id									SERIAL PRIMARY KEY,
		domain_id							INTEGER NOT NULL REFERENCES domains(domain_id) ON DELETE CASCADE,
		ip									VARCHAR(50) NOT NULL,
		created_at							TIMESTAMPTZ NOT NULL DEFAULT current_timestamp
	);

--####################INSERTING TEST DATA #####################
INSERT INTO domains(domain_name,forum_name,archived_at)VALUES('domain1','forum1',current_timestamp);
INSERT INTO domains(domain_name,forum_name,archived_at)VALUES('domain2','forum2',current_timestamp);
INSERT INTO domains(domain_name,forum_name,archived_at)VALUES('domain3','forum3',current_timestamp);

INSERT INTO users(domain_id,email,username,passwd_hash,onetime_login_token)VALUES(1,'user1@example.com','user1','pass1','token1');
INSERT INTO users(domain_id,email,username,passwd_hash,onetime_login_token)VALUES(2,'user2@example.com','user2','pass2','token2');
INSERT INTO users(domain_id,email,username,passwd_hash,onetime_login_token)VALUES(3,'user3@example.com','user3','pass3','token3');

INSERT INTO bannedips(domain_id,ip)VALUES(1,'192.168.1.1');
INSERT INTO bannedips(domain_id,ip)VALUES(2,'10.23.45.56');
INSERT INTO bannedips(domain_id,ip)VALUES(3,'134.45.67.187');