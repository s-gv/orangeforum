// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package models

import "github.com/jmoiron/sqlx"

func Migrate(db *sqlx.DB) error {
	migrateSqlite0001(db)
	return nil
}
