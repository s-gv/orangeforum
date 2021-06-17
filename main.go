// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
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

	createSuperUser := flag.Bool("createsuperuser", false, "Create superuser (interactive)")

	flag.Parse()

	db := sqlx.MustConnect("pgx", "postgres://dbuser:dbpass@localhost:5432/testdb")
	models.DB = db
	//models.Migrate()

	if *createSuperUser {
		fmt.Printf("Creating superuser...\n")
		userName, pass := getCreds()
		if userName != "" && pass != "" {
			if err := models.CreateSuperUser(userName, pass); err != nil {
				fmt.Printf("Error creating superuser: %s\n", err)
			}
		}
		return
	}

	views.SecretKey = os.Getenv("SECRET_KEY") // Ex: "s6JM1e8JTAphtKNR2y27XA8kkAaXOSYB" // 32 byte long
	if len(views.SecretKey) != 32 {
		glog.Errorf("Invalid Secret Key. Using randomly generated key. This will invalidate any active sessions.")

		var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
		b := make([]rune, 32)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		views.SecretKey = string(b)
	}

	r := views.GetRouter()

	glog.Info("Starting server on port 9123")
	http.ListenAndServe(":9123", r)
}
