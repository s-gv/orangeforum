package templates

const signupsrc = `
{{ define "content" }}

<form action="/signup" method="POST">
<input type="hidden" name="csrf" value="{{ .CSRF }}">
Username: <input type="text" name="username" required>
Password: <input type="password" name="passwd" required>
Confirm password: <input type="password" name="confirm" required>
Email (optional): <input type="text" name="email">
<input type="submit" value="Signup">
</form>

{{ .Msg }}

{{ end }}`