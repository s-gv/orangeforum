// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/views"
)

var TestServer *httptest.Server

func TestMain(m *testing.M) {
	models.DB = sqlx.MustConnect("pgx", "postgres://dbuser:dbpass@localhost:5432/oftestdb")
	views.SecretKey = "s6JM1e8JTAphtKNR2y27XA8kkAaXOSYB"

	router := views.GetRouter(true)
	TestServer = httptest.NewServer(router)
	defer TestServer.Close()

	exitVal := m.Run()
	os.Exit(exitVal)
}
