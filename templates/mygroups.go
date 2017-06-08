package templates

const mygroupsSrc = `
{{ define "content" }}
top | new | groups | create group
{{ if .IsUserValid }}
{{ .UserName }} ({{ .Karma }}) | <a href="/logout">logout</a>
{{ else }}
<a href="/login">login</a>
{{ end }}
{{ end }}`
