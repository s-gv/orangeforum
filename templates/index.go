package templates

const indexSrc = `
{{ define "content" }}


<h1>Groups</h1>
{{ range .Groups }}
<div class="row">
	<div><a href="/groups?name={{ .Name }}">{{ .Name }}</a></div>
	<div class="muted">{{ .Desc }}</div>
</div>
{{ end }}

<div class="row">
	{{ if not .GroupCreationDisabled }}
	<a class="link-btn" href="/groups/edit">New Group</a>
	{{ end }}
</div>

{{ end }}`