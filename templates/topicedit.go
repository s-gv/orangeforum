package templates

const topiceditSrc = `
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

{{ if not .TID }}
<h1>New topic</h1>
{{ else }}
<h1>Edit topic</h1>
{{ end }}


<form action="{{ if .TID }}/topics/edit{{ else }}/topics/new{{ end }}" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="id" value="{{ .TID }}">
<input type="hidden" name="gid" value="{{ .GID }}">

<div class="row clearfix">
	<div class="col1">Title</div>
	<div class="col2"><input type="text" name="title" placeholder="How does X work?" value="{{ .Title }}"></div>
</div>

<div class="row clearfix">
	<div class="col1">Content</div>
	<div class="col2"><textarea name="content" rows="4">{{ .Content }}</textarea></div>
</div>

{{ if .Common.Msg }}
<div class="row clearfix">
	<div class="col1-offset col2">
	{{ .Common.Msg }}
	</div>
</div>
{{ end }}

<div class="row clearfix">
	<div class="col1-offset col2">
	{{ if .TID }}
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
</div>

</form>

{{ end }}`
