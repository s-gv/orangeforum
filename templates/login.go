package templates

const loginsrc = `
{{ define "content" }}

<form action="/login" method="POST">
<input type="hidden" name="csrf" value="{{ .CSRF }}">
<input type="hidden" name="next" value="{{ .next }}">
Username: <input type="text" name="username">
Password: <input type="password" name="passwd">
<input type="submit" value="Login">
</form>
{{ .Msg }}

{{ end }}`