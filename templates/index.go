package templates

const indexsrc = `
{{ define "content" }}
Hello {{ .Name }}!
{{ .Msg }}

<form action="/test" method="POST">
<input type="submit" value="Test">
</form>

<a href="/signup">Signup</a>

<form action="/login" method="POST">
<input type="submit" value="Login">
</form>

<a href="/logout">Logout</a>.

{{ end }}`