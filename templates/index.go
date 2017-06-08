package templates

const indexSrc = `
{{ define "content" }}
<a href="/">top</a> | <a href="/new">new</a> | <a href="/groups">groups</a> | <a href="/creategroup">create group</a>
{{ if .IsUserValid }}
{{ .UserName }} ({{ .Karma }}) | <a href="/logout">logout</a>
{{ else }}
<a href="/login">login</a>
{{ end }}
{{ end }}`