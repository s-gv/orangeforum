package templates

const groupnewSrc = `
{{ define "content" }}
top | new | groups | submit
{{ if .IsUserValid }}
{{ .UserName }} ({{ .Karma }}) | <a href="/logout">logout</a>
{{ else }}
<a href="/login">login</a>
{{ end }}
{{ end }}`
