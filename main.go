// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/views"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	flag.Parse()

	db := sqlx.MustConnect("sqlite3", "orangeforum.db?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=-128000&_temp_store=2&_busy_timeout=2000")
	models.DB = db
	models.Migrate()

	views.SecretKey = os.Getenv("SECRET_KEY") // "s6JM1e8JTAphtKNR2y27XA8kkAaXOSYB" // 32 byte long

	r := views.GetRouter()

	glog.Info("Starting server on port 9123")
	http.ListenAndServe(":9123", r)
}
