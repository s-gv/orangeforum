package templates

const extranoteSrc = `
{{ define "content" }}
top | new | groups | create group
{{ if .IsUserValid }}
{{ .UserName }} ({{ .Karma }}) | <a href="/logout">logout</a>
{{ else }}
<a href="/login">login</a>
{{ end }}

<br>
{{ .ExtraNote.Name }}
<br>
{{ .ExtraNote.UpdatedDate }}
<br>
{{ .ExtraNote.Content }}
<br>

{{ end }}`