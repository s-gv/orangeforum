// Copyright (c) 2021 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed form.html
var formTmplStr string

type User struct {
	Name  string
	Email string
}

func ins(db *sqlx.DB) {
	name := "sagar"
	email := "sagar@example.com"

	tx, _ := db.Begin()
	for i := 0; i < 10; i++ {
		db.Exec("INSERT INTO users(name, email) VALUES($1, $2);", name, email)

		_, err := db.Exec("INSERT INTO users(name, email) VALUES($1, $2);", name, email)
		if err != nil {
			panic(err)
		}

		var users []User
		db.Select(&users, "SELECT name, email FROM users LIMIT 10;")
		println(len(users))
	}
	tx.Commit()
}

func main() {
	db := sqlx.MustConnect("sqlite3", "orangeforum.db?_journal_mode=WAL")

	db.MustExec(`CREATE TABLE users (
		name text,
		email text);`)

	ins(db)

	formTmpl := template.Must(template.New("form").Parse(formTmplStr))

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	id := uuid.New()
	fmt.Println(id.String())

	var tokenAuth *jwtauth.JWTAuth
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	// For debugging/example purposes, we generate and print
	// a sample jwt token with claims `user_id:123` here:
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123, "exp": time.Now()})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)

	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(tokenAuth))

		// Handle valid / invalid tokens. In this example, we use
		// the provided authenticator middleware, but you can write your
		// own very easily, look at the Authenticator method in jwtauth.go
		// and tweak it, its not scary.
		r.Use(jwtauth.Authenticator)

		r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			w.Write([]byte(fmt.Sprintf("protected area. hi %v", claims["user_id"])))
		})
	})

	csrfMiddleware := csrf.Protect([]byte("32-byte-long-auth-key"), csrf.Secure(false))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(csrfMiddleware)
		r.Use(sessionManager.LoadAndSave)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome anonymous"))
		})

		r.Get("/signin", func(w http.ResponseWriter, r *http.Request) {
			cookie := http.Cookie{Name: "jwt", Value: tokenString, Expires: time.Now().Add(365 * 24 * time.Hour)}
			http.SetCookie(w, &cookie)
			w.Write([]byte("signed in"))
		})

		r.Get("/form", func(w http.ResponseWriter, r *http.Request) {
			name := sessionManager.Pop(r.Context(), "name")
			formTmpl.Execute(w, map[string]interface{}{
				csrf.TemplateTag: csrf.TemplateField(r),
				"name":           name,
			})
		})

		r.Post("/form", func(w http.ResponseWriter, r *http.Request) {
			name := r.FormValue("name")
			sessionManager.Put(r.Context(), "name", name)
			http.Redirect(w, r, "/form", http.StatusSeeOther)
			// w.Write([]byte("hi " + name))
		})
	})

	addr := ":8000"
	fmt.Printf("Starting server on %v\n", addr)
	http.ListenAndServe(addr, r)
}
