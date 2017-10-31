// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const groupeditSrc = `
{{ define "content" }}

{{ if not .ID }}
<h1>New group</h1>
{{ else }}
<h1 id="title"><a href="/groups?name={{ .GroupName }}">{{ .GroupName }}</a></h1>
{{ end }}


<form action="/groups/edit" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="id" value="{{ .ID }}">
<table class="form">
	<tr>
		<th><label for="name">Group Name:</label></th>
		<td><input type="text" name="name" id="name" placeholder="Off-topic" value="{{ .GroupName }}"></td>
	</tr>
	<tr>
		<th><label for="desc">Description:</label></th>
		<td><input type="text" name="desc" id="desc" placeholder="Discuss topics not suitable in other groups" value="{{ .Desc }}"></td>
	</tr>
	<tr>
		<th><label for="header_msg">Announcement:</label></th>
		<td><textarea name="header_msg" id="header_msg" rows="4">{{ .HeaderMsg }}</textarea></td>
	</tr>
	<tr>
		<th><label for="mods">Mods:</label></th>
		<td><input type="text" name="mods" id="mods" placeholder="user1, user2" value="{{ .Mods }}"></td>
	</tr>
	<tr>
		<th><label for="admins">Admins (can edit this page):</label></th>
		<td><input type="text" name="admins" id="admins" placeholder="user1, user2" value="{{ .Admins }}"></td>
	</tr>
{{ if .Common.IsSuperAdmin }}
	<tr>
		<th><label for="is_sticky">Sticky:</label></th>
		<td><input type="checkbox" name="is_sticky" id="is_sticky"{{ if .IsSticky }} value="1" checked{{ end }}></td>
	</tr>
{{ end }}
{{ if .Common.Msg }}
	<tr>
		<th></th>
		<td>{{ .Common.Msg }}</td>
	</tr>
{{ end }}
	<tr>
		<th></th>
		<td>
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
		</td>
	</tr>
</table>
</form>

{{ end }}`
