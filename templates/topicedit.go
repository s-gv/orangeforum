// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const topiceditSrc = `
{{ define "content" }}

{{ if not .TopicID }}
<h1>New topic</h1>
{{ else }}
<h1 id="title"><a href="/topics?id={{ .TopicID }}">{{ .Title }}</a></h1>
{{ end }}


<form action="{{ if .TopicID }}/topics/edit{{ else }}/topics/new{{ end }}" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="id" value="{{ .TopicID }}">
<input type="hidden" name="gid" value="{{ .GroupID }}">
<table class="form">
	<tr>
		<th>Title:</th>
		<td><input type="text" name="title" placeholder="How does X work?" value="{{ .Title }}"></td>
	</tr>
	<tr>
		<th>Content:</th>
		<td><textarea name="content" rows="12">{{ .Content }}</textarea></td>
	</tr>
{{ if or .IsMod .IsAdmin .IsSuperAdmin }}
	<tr>
		<th>Sticky:</th>
		<td><input type="checkbox" name="is_sticky"{{ if .IsSticky }} checked{{ end }}></td>
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
		{{ if .TopicID }}
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
