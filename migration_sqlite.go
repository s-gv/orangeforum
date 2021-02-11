// Copyright (c) 2021 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import "github.com/jmoiron/sqlx"

func migrationSqlite0001(db *sqlx.DB) {
	db.MustExec(`CREATE TABLE users(
		id UUID PRIMARY KEY,
		email VARCHAR(250) DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT (STRFTIME('%Y-%m-%d %H:%M:%f', 'NOW')),
		updated_at DATETIME NOT NULL DEFAULT current_timestamp,
		archived_at DATETIME
	);`)
	db.MustExec(`CREATE TRIGGER updated_at_trig AFTER UPDATE ON users
		BEGIN UPDATE users SET updated_at = current_timestamp WHERE id = NEW.id; END;
	`)
}
