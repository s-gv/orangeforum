package templates

const indexSrc = `
{{ define "content" }}

<div class="btn-row">
	{{ if not .GroupCreationDisabled }}
	<a class="link-btn" href="/groups/edit">New Group</a>
	{{ end }}
</div>

<h1>Groups</h1>
<h2>{{ .HeaderMsg }}</h2>
{{ range .Groups }}
<div class="row">
	<div><a href="/groups?name={{ .Name }}">{{ .Name }}</a></div>
	<div class="muted">{{ .Desc }}</div>
</div>
{{ end }}

{{ end }}`