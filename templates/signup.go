// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const signupSrc = `
{{ define "content" }}

<form action="/signup" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="next" value="{{ .next }}">
<table class="form">
	<tr>
		<th><label for="username">Username:</label></th>
		<td><input type="text" name="username" id="username" required></td>
	</tr>
	<tr>
		<th><label for="passwd">Password:</label></th>
		<td><input type="password" name="passwd" id="passwd" required></td>
	</tr>
	<tr>
		<th><label for="confirm">Confirm password:</label></th>
		<td><input type="password" name="confirm" id="confirm" required></td>
	</tr>
	<tr>
		<th><label for="email">Email (optional):</label></th>
		<td><input type="text" name="email" id="email"></td>
	</tr>
	<tr>
		<th></th>
		<td>Already have an account? <a href="/login?next={{ .next }}">Login</a></td>
	</tr>
{{ if .Common.Msg }}
	<tr>
		<th></th>
		<td>{{ .Common.Msg }}</td>
	</tr>
{{ end }}
	<tr>
		<th></th>
		<td><input type="submit" value="Signup"></td>
	</tr>
</table>
</form>

{{ end }}`
