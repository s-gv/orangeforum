// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"github.com/s-gv/orangeforum/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := sqlx.MustConnect("sqlite3", "orangeforum.db?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=-128000&_temp_store=2&_busy_timeout=2000")
	models.Migrate(db)
}
