package templates

const indexSrc = `
{{ define "content" }}

{{ if not .GroupCreationDisabled }}
<a class="link-btn" style="float: right;" href="/groups/edit">New Group</a>
{{ end }}

{{ if .ShowGroups }}
<h1>Groups</h1>
{{ range .Groups }}
<div class="row">
	<div><a href="/groups?name={{ .Name }}">{{ .Name }}</a></div>
	<div class="muted">{{ .Desc }}</div>
</div>
{{ end }}
{{ end }}

{{ end }}`