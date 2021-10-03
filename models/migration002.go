// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import "github.com/jmoiron/sqlx"

func migrate002(db *sqlx.DB) {
	db.MustExec(`ALTER TABLE domains
		ADD COLUMN is_private                                      BOOL NOT NULL DEFAULT false,
		ADD COLUMN is_regular_signin_enabled                       BOOL NOT NULL DEFAULT true,
		ADD COLUMN is_auto_user_creation_on_email_signin_enabled   BOOL NOT NULL DEFAULT false,
		ADD COLUMN whitelisted_email_domains                       VARCHAR(250) NOT NULL DEFAULT '';
	`)
	db.MustExec(`UPDATE configs SET val = $1 WHERE name = $2;`, "2", DBVersion)
}
