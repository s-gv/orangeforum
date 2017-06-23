package templates

const profiletopicsSrc = `
{{ define "content" }}

<h2>Topics by {{ .OwnerName }}</h2>

{{ if .Topics }}
{{ range .Topics }}
<div class="row">
	<div>
		<a href="/topics?id={{ .ID }}">{{ .Title }}</a>{{ if .IsClosed }} [closed]{{ end }}
	</div>
	<div class="muted">{{ .CreatedDate }}</div>
</div>
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No topics to show.</div>
</div>
{{ end }}

{{ if .LastTopicDate }}
<div class="row">
	<div>
		<a href="/users/topics?u={{ .OwnerName }}&ltd={{ .LastTopicDate }}">More</a>
	</div>
</div>
{{ end }}

{{ end }}
`