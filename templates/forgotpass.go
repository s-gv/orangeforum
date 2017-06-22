package templates

const forgotpassSrc = `
{{ define "content" }}

<form action="/forgotpass" method="POST">
<input type="hidden" name="csrf" value={{ .Common.CSRF }}>
<table class="form">
	<tr>
		<th><label for="username">Username:</label></th>
		<td><input type="text" name="username" id="username" required></td>
	</tr>
{{ if .Common.Msg }}
	<tr>
		<th></th>
		<td>{{ .Common.Msg }}</td>
	</tr>
{{ end }}
	<tr>
		<th></th>
		<td><input type="submit" value="E-mail password reset link"></td>
	</tr>
</table>
</form>


{{ end }}`