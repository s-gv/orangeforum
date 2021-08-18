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
	ctxDomain   = contextKey("domain")
)

const BasePathField = "BasePath"
const DomainField = "Domain"

func domainCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainName := chi.URLParam(r, "domainName")
		domain := models.GetDomainByName(domainName)
		if domain == nil {
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			return
		}
		if domainName == r.Host {
			basePath := "/forums/" + domainName
			http.Redirect(w, r, r.URL.Path[len(basePath):], http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxDomain, domain)
		ctx2 := context.WithValue(ctx, ctxBasePath, "/forums/"+domainName)
		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}

func hostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainName := r.Host
		domain := models.GetDomainByName(domainName)
		if domain == nil {
			http.Error(w, http.StatusText(404), http.StatusNotFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxDomain, domain)
		ctx2 := context.WithValue(ctx, ctxBasePath, "")
		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}

func basePath(r *http.Request) string {
	return r.Context().Value(ctxBasePath).(string)
}
