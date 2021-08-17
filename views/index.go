// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/s-gv/orangeforum/models"
	"github.com/s-gv/orangeforum/templates"
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	domain := r.Context().Value(ctxDomain).(*models.Domain)
	user, _ := r.Context().Value(CtxUserKey).(*models.User)

	categories := models.GetCategoriesByDomainID(domain.DomainID)

	templates.Index.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath(r),
		DomainField:      domain,
		UserField:        user,
		"Categories":     categories,
	})
}
