// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const profileSrc = `
{{ define "content" }}


<form action="/users/update" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="u" value="{{ .UserName }}">
<table class="form">
	<tr>
		<th>Username:</th>
		<td>{{ .UserName }}</td>
	</tr>
	<tr>
		<th>About{{ if or .IsSelf .Common.IsSuperAdmin }} (public){{ end }}:</th>
		<td>{{ if or .IsSelf .Common.IsSuperAdmin }}<textarea name="about" rows="6">{{ .About }}</textarea>{{ else }}{{ .About }}{{ end }}
		</td>
	</tr>
{{ if or .IsSelf .Common.IsSuperAdmin }}
	<tr>
		<th><label for="email">Email (private):</label></th>
		<td><input type="email" name="email" id="email" value={{ .Email }}></td>
	</tr>
	{{ if .Common.Msg }}
	<tr>
		<th></th>
		<td><span class="alert">{{ .Common.Msg }}</span></td>
	</tr>
	{{ end }}
	<tr>
		<th></th>
		<td><input type="submit" name="action" value="Update"></td>
	</tr>
{{ end }}
	<tr>
		<th><a href="/users/topics?u={{ .UserName }}">topics</a>{{ if or .IsSelf .Common.IsSuperAdmin }} (public){{ end }}</th>
		<td></td>
	</tr>
	<tr>
		<th><a href="/users/comments?u={{ .UserName }}">comments</a>{{ if or .IsSelf .Common.IsSuperAdmin }} (public){{ end }}</th>
		<td></td>
	</tr>
{{ if or .IsSelf .Common.IsSuperAdmin }}
	<tr>
		<th><a href="/users/groups">groups</a>{{ if or .IsSelf .Common.IsSuperAdmin }} (private){{ end }}</th>
		<td></td>
	</tr>
	<tr>
		<th><a href="/changepass?u={{ .UserName }}">change password</a></th>
		<td></td>
	</tr>
{{ end }}
{{ if and .IsSelf .Common.IsSuperAdmin }}
	<tr>
		<th><a href="/admin">admin section</a></th>
		<td></td>
	</tr>
	<tr>
		<th><a href="/signup">sign-up new user</a>
		<td></td>
	</tr>
{{ end }}
{{ if .IsSelf }}
	<tr>
		<th><a href="/pm">private messages</a></th>
		<td></td>
	</tr>
	<tr>
		<th><a href="/logout">logout</a></th>
		<td></td>
	</tr>
{{ end }}
{{ if .Common.IsSuperAdmin }}
{{ if not .IsSelf }}
	<tr>
		<th></th>
		<td>
			{{ if not .IsBanned }}
			<input type="submit" name="action" value="Ban">
			{{ else }}
			<input type="submit" name="action" value="Unban">
			{{ end }}
		</td>
	</tr>
{{ end }}
{{ end }}
</table>
</form>

{{ end }}`
