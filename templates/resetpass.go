// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const resetpassSrc = `
{{ define "content" }}

<form action="/resetpass" method="POST">
<input type="hidden" name="csrf" value={{ .Common.CSRF }}>
<input type="hidden" name="r" value={{ .ResetToken }}>
<table class="form">
	<tr>
		<th><label for="passwd">Password:</label></th>
		<td><input type="password" name="passwd" id="passwd" required></td>
	</tr>
	<tr>
		<th><label for="confirm">Confirm password:</label></th>
		<td><input type="password" name="confirm" id="confirm" required></td>
	</tr>
{{ if .Common.Msg }}
	<tr>
		<th></th>
		<td>{{ .Common.Msg }}</td>
	</tr>
{{ end }}
	<tr>
		<th></th>
		<td><input type="submit" value="Change password"></td>
	</tr>
</table>
</form>

{{ end }}`
