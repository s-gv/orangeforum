package templates

const profilegroupsSrc = `
{{ define "content" }}


<h1>{{ .OwnerName }}</h1>

<h2>Admin</h2>
{{ if .AdminInGroups }}
{{ range .AdminInGroups }}
<div class="row">
	<div>
		<a href="/groups?id={{ .ID }}">{{ .Name }}</a>{{ if .IsClosed }} [closed]{{ end }}
	</div>
	<div class="muted">{{ .CreatedDate }}</div>
</div>
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No groups to show.</div>
</div>
{{ end }}

<h2>Mod</h2>
{{ if .ModInGroups }}
{{ range .ModInGroups }}
<div class="row">
	<div>
		<a href="/groups?id={{ .ID }}">{{ .Name }}</a>{{ if .IsClosed }} [closed]{{ end }}
	</div>
	<div class="muted">{{ .CreatedDate }}</div>
</div>
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No groups to show.</div>
</div>
{{ end }}

{{ end }}
`