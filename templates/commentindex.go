package templates

const commentindexSrc = `
{{ define "content" }}

<h1><a href="/groups?name={{ .GroupName }}">{{ .GroupName }}</a></h1>
<h2><a href="/topics?id={{ .TopicID }}">{{ .TopicName }}</a></h2>

<div class="row">
	<div>by {{ .UserName }} <a href="/comments?id={{ .ID }}">{{ .CreatedDate }}</a>{{ if or .IsOwner .IsAdmin .IsMod .IsSuperAdmin }} | <a href="/comments/edit?id={{ .ID }}">edit</a> {{end}}</div>
	{{ if .IsDeleted }}
		<div>[DELETED]</div>
	{{ else }}
		<div>{{ .Content }}</div>
		{{ if .ImgSrc }}<div><img src="/img?name={{ .ImgSrc }}"></div>{{ end }}
	{{ end }}
</div>


{{ end }}`
