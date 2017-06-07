package templates

const indexsrc = `
{{ define "content" }}
Hello {{ .Name }}!
{{ .Msg }}

<form action="/test" method="POST">
<input type="submit" value="Test">
</form>

<a href="/login">Login</a>
<a href="/signup">Signup</a>
<a href="/logout">Logout</a>.

{{ end }}`