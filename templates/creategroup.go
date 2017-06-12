package templates

const creategroupSrc = `
{{ define "content" }}
top | new | groups | create group
{{ if .IsUserValid }}
{{ .UserName }} ({{ .Karma }}) | <a href="/logout">logout</a>
{{ else }}
<a href="/login">login</a>
{{ end }}


<form action="/creategroup" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
Group Name: <input type="text" name="name">
Group Description: <input type="text" name="desc">
<input type="submit" value="Create Group">
</form>

{{ .Msg }}

{{ end }}`