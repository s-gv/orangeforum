// Copyright (c) 2021 Orange Forum authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package views

import (
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/s-gv/orangeforum/templates"
)

func getIndex(w http.ResponseWriter, r *http.Request) {
	templates.Index.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		BasePathField:    basePath(r),
	})
}
