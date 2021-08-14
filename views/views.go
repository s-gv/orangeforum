// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/csrf"
	"github.com/s-gv/orangeforum/models"
)

var tokenAuth *jwtauth.JWTAuth

// SecretKey must be 32 byte long.
var SecretKey string

func forumRouter() *chi.Mux {
	r := chi.NewRouter()

	csrfMiddleware := csrf.Protect([]byte(SecretKey), csrf.Secure(false))
	r.Use(csrfMiddleware)

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	r.Use(sessionManager.LoadAndSave)

	tokenAuth = jwtauth.New("HS256", []byte(SecretKey), nil)

	// Auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Route("/signin", func(r chi.Router) {
			r.Get("/", getAuthSignIn)
			r.Post("/", postAuthSignIn)
		})

		r.Route("/otsignin", func(r chi.Router) {
			r.Get("/", getAuthOneTimeSignIn)
			r.Post("/", postAuthOneTimeSignIn)
			r.Get("/{token}", getAuthOneTimeSignInDone)
		})

		r.Route("/signup", func(r chi.Router) {
			r.Get("/", getAuthSignUp)
			r.Post("/", postAuthSignUp)
		})

		r.Route("/signup/{signupToken:[A-Za-z0-9-]+}", func(r chi.Router) {
			r.Get("/", getAuthSignUp)
			r.Post("/", postAuthSignUp)
		})

		r.Route("/changepass", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(mustAuth)

			r.Get("/", getAuthChangePass)
			r.Post("/", postAuthChangePass)
		})

		r.Route("/logout", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(canAuth)

			r.Get("/", getAuthLogout)
		})
	})

	r.Route("/can", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(canAuth)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var user *models.User
			if u, ok := r.Context().Value(CtxUserKey).(*models.User); ok {
				user = u
			}
			domain, _ := r.Context().Value(ctxDomain).(*models.Domain)
			w.Write([]byte(fmt.Sprintf("public area. hi %v. domain_id: %v, domain name: %v", user, domain.DomainID, domain.DomainName)))
		})
	})

	r.Route("/must", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(mustAuth)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := r.Context().Value(CtxUserKey).(*models.User)
			domain, _ := r.Context().Value(ctxDomain).(*models.Domain)
			w.Write([]byte(fmt.Sprintf("protected area. \nuser: %v\ndomain: %v", user, domain)))
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(canAuth)

		r.Get("/", getIndex)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(mustAuth)

		r.Get("/", getAdmin)
		r.Post("/", postAdmin)

		r.Post("/mods/create", postCreateMod)
		r.Post("/mods/delete", postDeleteMod)

		r.Post("/categories/{categoryID}", postCategory)
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(canAuth)

		r.Get("/{userID}", getProfile)
		r.Post("/{userID}", postProfile)
	})

	return r
}

func GetRouter() *chi.Mux {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	fr := forumRouter()

	r.Route("/domains/{domainName}", func(r chi.Router) {
		r.Use(domainCtx)
		r.Mount("/", fr)
	})

	r.Route("/", func(r chi.Router) {
		r.Use(hostCtx)
		r.Mount("/", fr)
	})

	return r
}
