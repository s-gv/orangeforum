// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const indexSrc = `
{{ define "content" }}

<div class="btn-row">
	{{ if not .GroupCreationDisabled }}
	<a class="link-btn" href="/groups/edit">New Group</a>
	{{ end }}
</div>

<h1>Groups</h1>
<h2>{{ .HeaderMsg }}</h2>
{{ if .Groups }}
{{ range .Groups }}
<div class="topic-row">
	<div><a href="/groups?name={{ .Name }}">{{ .Name }}</a></div>
	<div class="muted">{{ .Desc }}</div>
</div>
<hr class="sep">
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No groups to show.</div>
</div>
{{ end }}

<h1>Recent Topics</h1>
{{ if .Topics }}
{{ range .Topics }}
<div class="topic-row">
	<div><a href="/topics?id={{ .ID }}">{{ .Title }}</a></div>
	<div class="muted">
		<a href="/users?u={{ .OwnerName }}">{{ .OwnerName }}</a> in <a href="/groups?name={{ .GroupName }}">{{ .GroupName }}</a> {{ .CreatedDate }} | <a href="/topics?id={{ .ID }}">{{ .NumComments }} comments</a>
	</div>
</div>
<hr class="sep">
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No topics to show.</div>
</div>
{{ end }}

{{ end }}`
