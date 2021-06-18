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
	BasePath = contextKey("base_path")
	DomainID = contextKey("domain_id")
)

func DomainCtx(next http.Handler) http.Handler {
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
		ctx := context.WithValue(r.Context(), DomainID, *domainID)
		ctx2 := context.WithValue(ctx, BasePath, "/domains/"+domainName)
		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}

func HostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainName := r.Host
		domainID := models.GetDomainIDByName(domainName)
		if domainID == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		ctx := context.WithValue(r.Context(), DomainID, *domainID)
		ctx2 := context.WithValue(ctx, BasePath, "")
		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}
