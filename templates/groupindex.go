package templates

const groupindexSrc = `
{{ define "content" }}

<h1><a href="/groups?name={{ .GroupName }}">{{ .GroupName }}</a></h1>
{{ if or .IsAdmin .IsMod .IsSuperAdmin }}
<a href="/groups/edit?id={{ .GroupID }}">edit</a>
{{ end }}

<a href="/topics/new?gid={{ .GroupID }}">New topic</a>

{{ if .Topics }}
{{ range .Topics }}
<div class="row">
	<div><a href="/topics?id={{ .ID }}">{{ .Title }}</a></div>
	<div class="muted">by {{ .Owner }} {{ .CreatedDate }} | {{ .NumComments }} comments</div>
</div>
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No topics here.</div>
</div>
{{ end }}

{{ if .LastTopicDate }}
<div class="row">
	<div><a href="/groups?name={{ .GroupName }}&ltd={{ .LastTopicDate }}">More</a></div>
</div>
{{ end }}

{{ end }}`
