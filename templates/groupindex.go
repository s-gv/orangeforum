package templates

const groupindexSrc = `
{{ define "content" }}

<h1>{{ .Name }}</h1>
{{ if or .IsAdmin .IsMod .IsSuperAdmin }}
<a href="/groups/edit?id={{ .GroupID }}">edit</a>
{{ end }}

<a href="/topics/new?gid={{ .GroupID }}">New topic</a>

{{ range .Topics }}
<div class="row">
	<div><a href="/topics?id={{ .ID }}">{{ .Title }}</a></div>
	<div class="muted">by {{ .Owner }} {{ .CreatedDate }} | {{ .NumComments }} comments</div>
</div>
{{ end }}

{{ end }}`
