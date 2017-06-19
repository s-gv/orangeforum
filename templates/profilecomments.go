package templates

const profilecommentsSrc = `
{{ define "content" }}

<h2>{{ .OwnerName }}</h2>

{{ if .Comments }}
{{ range .Comments }}
<div class="row">
	<div>
		by {{ $.OwnerName }} <a href="/comments?id={{ .ID }}">{{ .CreatedDate }}</a> on <a href="/topics?id={{ .TopicID }}">{{ .TopicName }}</a>
	</div>
	<div>{{ if .IsDeleted }}[DELETED]{{ else }}{{ .Content }}{{ end }}</div>
</div>
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No comments to show.</div>
</div>
{{ end }}

{{ if .LastCommentDate }}
<div class="row">
	<div>
		<a href="/users/comments?u={{ .OwnerName }}&lcd={{ .LastCommentDate }}">More</a>
	</div>
</div>
{{ end }}

{{ end }}
`