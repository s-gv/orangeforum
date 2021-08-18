// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/gorilla/csrf"
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

	r.Route("/static", func(r chi.Router) {
		r.Get("/orangeforum.css", getCSS)
		r.Get("/orangeforum.js", getJS)
	})

	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(canAuth)

		r.Get("/", getIndex)
	})

	r.Route("/admin", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(mustAuth)

		r.Get("/", adminHandler)
		r.Post("/", adminHandler)

		r.Post("/mods/create", postCreateMod)
		r.Post("/mods/delete", postDeleteMod)

		r.Post("/categories/{categoryID}", postCategory)
	})

	r.Route("/users", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(canAuth)

		r.Get("/{userID}", profileHandler)
		r.Post("/{userID}", profileHandler)
	})

	r.Route("/categories/{categoryID}", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(canAuth)

			r.Get("/", getTopicList)
			r.Get("/topics/{topicID}", getTopic)
		})

		r.Route("/topics", func(r chi.Router) {
			r.Route("/new", func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(mustAuth)

				r.Get("/", editTopic)
				r.Post("/", editTopic)
			})

			r.Route("/{topicID}/edit", func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(mustAuth)

				r.Get("/", editTopic)
				r.Post("/", editTopic)
			})

			r.Route("/{topicID}/comments/new", func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(mustAuth)

				r.Get("/", editComment)
				r.Post("/", editComment)
			})

			r.Route("/{topicID}/comments/{commentID}/edit", func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(mustAuth)

				r.Get("/", editComment)
				r.Post("/", editComment)
			})

			r.Route("/{topicID}", func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Use(canAuth)

				r.Get("/", getTopic)
			})
		})

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

	r.Route("/forums/{domainName}", func(r chi.Router) {
		r.Use(domainCtx)
		r.Mount("/", fr)
	})

	r.Route("/", func(r chi.Router) {
		r.Use(hostCtx)
		r.Mount("/", fr)
	})

	return r
}
