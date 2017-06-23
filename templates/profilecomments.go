package templates

const profilecommentsSrc = `
{{ define "content" }}

<h2>Comments by {{ .OwnerName }}</h2>

{{ if .Comments }}
{{ range .Comments }}
<div class="row">
	<div class="muted">by {{ $.OwnerName }}</a> <a href="/comments?id={{ .ID }}">{{ .CreatedDate }}</a> on <a href="/topics?id={{ .TopicID }}">{{ .TopicName }}</a></div>
	{{ if .IsDeleted }}
		<div>[DELETED]</div>
	{{ else }}
		<div>{{ .Content }}</div>
		{{ if .ImgSrc }}<div><img src="/img?name={{ .ImgSrc }}"></div>{{ end }}
	{{ end }}
	<hr class="sep">
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