package templates

const indexsrc = `
{{ define "content" }}
Hello {{ .Name }}!
{{ .Msg }}

<form action="/signup" method="POST">
<input type="submit" value="Signup">
</form>

<form action="/login" method="POST">
<input type="submit" value="Login">
</form>

{{ end }}`