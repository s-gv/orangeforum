package templates

const indexsrc = `
{{ define "content" }}
Hello {{ .Name }}!
{{ .Msg }}

<form action="/test" method="POST">
<input type="submit" value="Test">
</form>

<form action="/signup" method="POST">
<input type="submit" value="Signup">
</form>

<form action="/login" method="POST">
<input type="submit" value="Login">
</form>

{{ end }}`