package templates

const resetpassSrc = `
{{ define "content" }}

<form action="/resetpass" method="POST">
<input type="hidden" name="csrf" value={{ .CSRF }}>
<input type="hidden" name="r" value={{ .ResetToken }}>
Password: <input type="password" name="passwd" required>
Confirm password: <input type="password" name="confirm" required>
<input type="submit" value="Change password">
</form>
{{ .Msg }}

{{ end }}`