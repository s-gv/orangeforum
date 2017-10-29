// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const extranoteSrc = `
{{ define "content" }}


<h1>{{ .Name }}</h1>
<p>{{ .UpdatedDate.Format "02 Jan 2006" }}</p>
<div>{{ .Content }}</div>

{{ end }}`
