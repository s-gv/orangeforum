package templates

const indexSrc = `
{{ define "content" }}

{{ if not .GroupCreationDisabled }}
<a href="/groups/edit">create group</a>
{{ end }}

{{ end }}`