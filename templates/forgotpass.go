package templates

const forgotpasssrc = `
{{ define "content" }}

<form action="/forgotpass" method="POST">
<input type="hidden" name="csrf" value={{ .CSRF }}>
Username: <input type="text" name="username" required>
<input type="submit" value="E-mail reset link">
</form>
{{ .Msg }}

{{ end }}`