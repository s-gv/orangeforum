package templates

const groupnewSrc = `
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

{{ if not .ID }}
<h1>New group</h1>
{{ else }}
<h1>Edit group</h1>
{{ end }}


<form action="/groups/edit" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="id" value="{{ .ID }}">

<div class="row clearfix">
	<div class="col1">Group Name</div>
	<div class="col2"><input type="text" name="name" placeholder="Off-topic" value="{{ .Name }}"></div>
</div>

<div class="row clearfix">
	<div class="col1">Description</div>
	<div class="col2"><input type="text" name="desc" placeholder="Discuss topics not suitable in other groups" value="{{ .Desc }}"></div>
</div>

<div class="row clearfix">
	<div class="col1">Announcement</div>
	<div class="col2"><textarea name="header_msg" rows="4">{{ .HeaderMsg }}</textarea></div>
</div>

<div class="row clearfix">
	<div class="col1">Mods</div>
	<div class="col2"><input type="text" name="mods" placeholder="user1, user2" value="{{ .Mods }}"></div>
</div>

<div class="row clearfix">
	<div class="col1">Admins (can edit this page)</div>
	<div class="col2"><input type="text" name="admins" placeholder="user1, user2" value="{{ .Admins }}"></div>
</div>
{{ if .Common.IsSuperAdmin }}
<div class="row clearfix">
	<div class="col1">Sticky</div>
	<div class="col2"><input type="checkbox" name="is_sticky"{{ if .IsSticky }} value="1" checked{{ end }}></div>
</div>
{{ end }}

{{ if .Common.Msg }}
<div class="row clearfix">
	<div class="col1-offset col2">
	{{ .Common.Msg }}
	</div>
</div>
{{ end }}

<div class="row clearfix">
	<div class="col1-offset col2">
	{{ if .ID }}
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
