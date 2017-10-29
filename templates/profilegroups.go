// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const profilegroupsSrc = `
{{ define "content" }}


<h1>{{ .OwnerName }}</h1>

<h2>Admin</h2>
{{ if .AdminInGroups }}
{{ range .AdminInGroups }}
<div class="row">
	<div>
		<a href="{{ if .IsClosed }}/groups/edit?id={{ .ID }}{{ else }}/groups?name={{ .Name }}{{ end }}">{{ .Name }}</a>{{ if .IsClosed }} [closed]{{ end }}
	</div>
	<div class="muted">created {{ .CreatedDate }}</div>
</div>
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No groups to show.</div>
</div>
{{ end }}

<h2>Mod</h2>
{{ if .ModInGroups }}
{{ range .ModInGroups }}
<div class="row">
	<div>
		<a href="{{ if .IsClosed }}/groups/edit?id={{ .ID }}{{ else }}/groups?name={{ .Name }}{{ end }}">{{ .Name }}</a>{{ if .IsClosed }} [closed]{{ end }}
	</div>
	<div class="muted">created {{ .CreatedDate }}</div>
</div>
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No groups to show.</div>
</div>
{{ end }}

{{ end }}
`
