package templates


const commenteditSrc = `
{{ define "head" }}
<style>
.row {
	margin-top: 15px;
}
input[type="text"], input[type="number"], textarea {
	width: 90%;
}
@media screen and (min-width:600px) {
	.col1 {
		float: left;
		text-align: right;
		width: 275px;
	}
	.col2 {
		float: left;
		text-align: left;
		margin-left: 15px;
		width: 300px;
	}
	.col1-offset {
		margin-left: 290px;
	}
}
</style>
{{ end }}


{{ define "content" }}

<h1>{{ .GroupName }}</h1>
<h2>{{ .TopicName }}</h2>

<div>{{ .ParentComment }}</div>

<form action="{{ if .CommentID }}/comments/edit{{ else }}/comments/new{{ end }}" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="id" value="{{ .CommentID }}">
<input type="hidden" name="tid" value="{{ .TopicID }}">
<div>
<textarea name="content" rows="8">{{ .Content }}</textarea>
</div>

{{ if or .IsMod .IsAdmin .IsSuperAdmin }}
<div><input type="checkbox" name="is_sticky"{{ if .IsSticky }} checked{{ end }}> Sticky</div>
{{ end }}

{{ .Common.Msg }}

<div>
{{ if .CommentID }}
	{{ if not .IsDeleted }}
	<input type="submit" name="action" value="Update">
	<input type="submit" name="action" value="Delete">
	{{ else }}
	<input type="submit" name="action" value="Undelete">
	{{ end }}
{{ else }}
<input type="submit" name="action" value="Create">
{{ end }}
</div>

</form>

{{ end }}`
