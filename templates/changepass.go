package templates

const changepassSrc = `
{{ define "content" }}

<form action="/changepass" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<table class="form">
	<tr>
		<th><label for="passwd">Current password:</label></th>
		<td><input type="password" name="passwd" id="passwd" required></td>
	</tr>
	<tr>
		<th><label for="newpass">New password:</label></th>
		<td><input type="password" name="newpass" id="newpass" required></td>
	</tr>
	<tr>
		<th><label for="confirm">New password (again):</label></th>
		<td><input type="password" name="confirm" id="confirm" required></td>
	</tr>
	<tr>
		<th></th>
		<td><input type="submit" value="Change Password"></td>
	</tr>
</table>
</form>


{{ .Msg }}

{{ end }}`

