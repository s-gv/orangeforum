// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const changepassSrc = `
{{ define "content" }}

<form action="/changepass" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="u" value="{{ .UserName }}">
<table class="form">
	<tr>
		<th><label for="username">User:</label></th>
		<td>{{ .UserName }}</td>
	</tr>
{{ if not .Common.IsSuperAdmin }}
	<tr>
		<th><label for="passwd">Current password:</label></th>
		<td><input type="password" name="passwd" id="passwd" required></td>
	</tr>
{{ end }}
	<tr>
		<th><label for="newpass">New password:</label></th>
		<td><input type="password" name="newpass" id="newpass" required></td>
	</tr>
	<tr>
		<th><label for="confirm">New password (again):</label></th>
		<td><input type="password" name="confirm" id="confirm" required></td>
	</tr>
{{ if .Common.Msg }}
	<tr>
		<th></th>
		<td><span class="alert">{{ .Common.Msg }}</span></td>
	</tr>
{{ end }}
	<tr>
		<th></th>
		<td><input type="submit" value="Change Password"></td>
	</tr>
</table>
</form>

{{ end }}`
