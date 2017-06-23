package templates

const commentindexSrc = `
{{ define "content" }}

<div class="row">
	<div class="muted">
		comment by <a href="/users?u={{ .OwnerName }}">{{ .OwnerName }}</a> to <a href="/topics?id={{ .TopicID }}">{{ .TopicName }}</a>
		<a href="/comments?id={{ .ID }}">{{ .CreatedDate }}</a>
		{{ if or .IsOwner $.IsAdmin $.IsMod $.IsSuperAdmin }} | <a href="/comments/edit?id={{ .ID }}">edit</a> {{end}}
	</div>
	{{ if .IsDeleted }}
		<div>[DELETED]</div>
	{{ else }}
		<div>{{ .Content }}</div>
		{{ if .ImgSrc }}<div><img src="/img?name={{ .ImgSrc }}"></div>{{ end }}
	{{ end }}
</div>


{{ end }}`
