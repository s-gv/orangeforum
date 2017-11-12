// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const loginSrc = `
{{ define "content" }}

<form action="/login" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="next" value="{{ .next }}">
<table class="form">
{{ if .LoginMsg }}
	<tr>
		<th></th>
		<td>{{ .LoginMsg }}</td>
	</tr>
{{ end }}
	<tr>
		<th>Username:</th>
		<td><input type="text" name="username" required></td>
	</tr>
	<tr>
		<th>Password:</th>
		<td><input type="password" name="passwd" required></td>
	</tr>
	<tr>
		<th></th>
		<td>Don't have an account? <a href="/signup?next={{ .next }}">Signup</a></td>
	</tr>
	<tr>
		<th></th>
		<td><a href="/forgotpass">Forgot password?</a></td>
	</tr>
{{ if .Common.Msg }}
	<tr>
		<th></th>
		<td><span class="alert">{{ .Common.Msg }}</span></td>
	</tr>
{{ end }}
	<tr>
		<th></th>
		<td><input type="submit" value="Login"></td>
	</tr>
</table>
</form>

{{ end }}`
