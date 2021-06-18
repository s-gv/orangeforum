// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"flag"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/views"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	secretKey := os.Getenv("SECRET_KEY") // Ex: "s6JM1e8JTAphtKNR2y27XA8kkAaXOSYB" // 32 byte long
	dsn := os.Getenv("ORANGEFORUM_DSN")

	port := flag.String("port", "9123", "Port to listen on")
	shouldMigrate := flag.Bool("migrate", false, "Migrate DB")
	createSuperUser := flag.Bool("createsuperuser", false, "Create superuser (interactive)")
	createUser := flag.Bool("createuser", false, "Create user. Optional arguments: <domain> <username> <password> <email>")
	changePasswd := flag.Bool("changepasswd", false, "Change password")

	flag.Parse()

	if dsn == "" {
		dsn = "postgres://dbuser:dbpass@localhost:5432/testdb"
		glog.Infof("Environment variable ORANGEFORUM_DSN not set. Using default dsn: %s", dsn)
	}

	db := sqlx.MustConnect("pgx", dsn)
	models.DB = db

	if *shouldMigrate {
		err := models.Migrate()
		if err != nil {
			glog.Error(err)
		}
		return
	}

	err := models.IsMigrationNeeded()
	if err != nil {
		glog.Error(err)
		return
	}

	if *createSuperUser {
		commandCreateSuperUser()
		return
	}

	if *changePasswd {
		commandChangePasswd()
		return
	}

	if *createUser {
		commandCreateUser()
		return
	}

	if len(secretKey) != 32 {
		glog.Errorf("Secret key in environment variable SECRET_KEY does not have length 32. Using randomly generated key. This will invalidate any active sessions.")

		var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		b := make([]rune, 32)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		secretKey = string(b)
	}

	views.SecretKey = secretKey
	r := views.GetRouter()

	glog.Info("Starting server on port " + *port)
	http.ListenAndServe(":"+*port, r)
}
