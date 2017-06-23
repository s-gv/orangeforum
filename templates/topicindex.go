package templates

const topicindexSrc = `
{{ define "content" }}

<div class="btn-row">
	{{ if not .IsClosed }}
	<a class="link-btn" href="/comments/new?tid={{ .TopicID }}">Reply</a>
	{{ end }}
	{{ if or .IsAdmin .IsMod .IsSuperAdmin .IsOwner }}
	<a class="link-btn" href="/topics/edit?id={{ .TopicID }}">Edit topic</a>
	{{ end }}
	{{ if and .Common.UserName .Common.IsTopicSubAllowed }}
	{{ if .SubToken }}
	<form action="/topics/unsubscribe?token={{ .SubToken }}" method="POST">
		<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
		<input class="btn" type="submit" value="Unsubscribe">
	</form>
	{{ else }}
	<form action="/topics/subscribe?id={{ .TopicID }}" method="POST">
		<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
		<input class="btn" type="submit" value="Subscribe">
	</form>
	{{ end }}
	{{ end }}
</div>

<h1><a href="/groups?name={{ .GroupName }}">{{ .GroupName }}</a></h1>
<h2>{{ .TopicName }}</h2>

<div>
{{ .Content }}
</div>







{{ if .Comments }}
{{ range .Comments }}
<div class="row">
	<div>by {{ .UserName }} <a href="/comments?id={{ .ID }}">{{ .CreatedDate }}</a>{{ if or .IsOwner $.IsAdmin $.IsMod $.IsSuperAdmin }} | <a href="/comments/edit?id={{ .ID }}">edit</a> {{end}}</div>
	{{ if .IsDeleted }}
		<div>[DELETED]</div>
	{{ else }}
		<div>{{ .Content }}</div>
		{{ if .ImgSrc }}<div><img src="/img?name={{ .ImgSrc }}"></div>{{ end }}
	{{ end }}
</div>
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No comments to show.</div>
</div>
{{ end }}

{{ if .LastCommentDate }}
<div class="row">
	<div><a href="/topics?id={{ .TopicID }}&lcd={{ .LastCommentDate }}">More</a></div>
</div>
{{ end }}

{{ end }}`
