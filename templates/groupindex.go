package templates

const groupindexSrc = `
{{ define "content" }}

<div class="btn-row">
	<a class="link-btn" href="/topics/new?gid={{ .GroupID }}">New topic</a>
	{{ if or .IsAdmin .IsMod .IsSuperAdmin }}
	<a class="link-btn" href="/groups/edit?id={{ .GroupID }}">Edit group</a>
	{{ end }}
	{{ if and .Common.UserName .Common.IsGroupSubAllowed }}
	{{ if .SubToken }}
	<form action="/groups/unsubscribe?token={{ .SubToken }}" method="POST">
		<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
		<input class="btn" type="submit" value="Unsubscribe">
	</form>
	{{ else }}
	<form action="/groups/subscribe?id={{ .GroupID }}" method="POST">
		<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
		<input class="btn" type="submit" value="Subscribe">
	</form>
	{{ end }}
	{{ end }}
</div>

<h1>{{ .GroupName }}</h1>
<h2>{{ .HeaderMsg }}</h2>

{{ if .Topics }}
{{ range .Topics }}
<div class="row">
	<div><a href="/topics?id={{ .ID }}">{{ .Title }}</a></div>
	<div class="muted">by {{ .Owner }} {{ .CreatedDate }} | {{ .NumComments }} comments</div>
</div>
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No topics here.</div>
</div>
{{ end }}

{{ if .LastTopicDate }}
<div class="row">
	<div><a href="/groups?name={{ .GroupName }}&ltd={{ .LastTopicDate }}">More</a></div>
</div>
{{ end }}

{{ end }}`
