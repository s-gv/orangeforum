package templates

const changepasssrc = `
{{ define "content" }}

<form action="/changepass" method="POST">
<input type="hidden" name="csrf" value="{{ .CSRF }}">
Current password: <input type="password" name="passwd" required>
New password: <input type="password" name="newpass" required>
New password (again): <input type="password" name="confirm" required>
<input type="submit" value="Change Password">
</form>

{{ .Msg }}

{{ end }}`

