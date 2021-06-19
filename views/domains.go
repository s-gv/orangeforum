// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/s-gv/orangeforum/models"
)

type contextKey string

const (
	ctxBasePath = contextKey("base_path")
	ctxDomainID = contextKey("domain_id")
)

const BasePathField = "BasePath"

func domainCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainName := chi.URLParam(r, "domainName")
		domainID := models.GetDomainIDByName(domainName)
		if domainID == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		if domainName == r.Host {
			basePath := "/domains/" + domainName
			http.Redirect(w, r, r.URL.Path[len(basePath):], http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxDomainID, *domainID)
		ctx2 := context.WithValue(ctx, ctxBasePath, "/domains/"+domainName)
		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}

func hostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainName := r.Host
		domainID := models.GetDomainIDByName(domainName)
		if domainID == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), ctxDomainID, *domainID)
		ctx2 := context.WithValue(ctx, ctxBasePath, "")
		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}

func basePath(r *http.Request) string {
	return r.Context().Value(ctxBasePath).(string)
}

func domainID(r *http.Request) int {
	return r.Context().Value(ctxDomainID).(int)
}
