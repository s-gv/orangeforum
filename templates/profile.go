package templates

const profileSrc = `
{{ define "content" }}


<form action="/users" method="POST">
<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
<input type="hidden" name="u" value="{{ .UserName }}">
<table class="form">
	<tr>
		<th>Username:</th>
		<td>{{ .UserName }}</td>
	</tr>
	<tr>
		<th>About{{ if .IsSelf }} (public):{{ end }}</th>
		<td>{{ if .IsSelf }}<textarea name="about" rows="6">{{ .About }}</textarea>{{ else }}{{ .About }}{{ end }}
		</td>
	</tr>
{{ if .IsSelf }}
	<tr>
		<th><label for="email">Email (private):</label></th>
		<td><input type="email" name="email" id="email" value={{ .Email }}></td>
	</tr>
	<tr>
		<th></th>
		<td><input type="submit" name="action" value="Update"></td>
	</tr>
{{ end }}
	<tr>
		<th><a href="/users/topics?u={{ .UserName }}">topics</a>{{ if .IsSelf }} (public){{ end }}</th>
		<td></td>
	</tr>
	<tr>
		<th><a href="/users/comments?u={{ .UserName }}">comments</a>{{ if .IsSelf }} (public){{ end }}</th>
		<td></td>
	</tr>
{{ if .IsSelf }}
	<tr>
		<th><a href="/users/groups">groups</a>{{ if .IsSelf }} (private){{ end }}</th>
		<td></td>
	</tr>
	<tr>
		<th><a href="/changepass">change password</a></th>
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
{{ if and .IsSelf .Common.IsSuperAdmin }}
	<tr>
		<th><a href="/admin">admin</a> (private)</th>
		<td></td>
	</tr>
{{ end }}
</table>
</form>

{{ end }}`
