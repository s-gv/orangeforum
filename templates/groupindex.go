package templates

const groupindexSrc = `
{{ define "content" }}

<h1>{{ .Name }}</h1>
{{ if or .IsAdmin .Common.IsSuperAdmin }}
<a href="/groups/edit?id={{ .ID }}">edit</a>
{{ end }}

{{ end }}`
