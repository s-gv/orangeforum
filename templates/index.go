package templates

const indexsrc = `
{{ define "content" }}
Hello world!
{{ .Msg }}
<form action="/test" method="POST">
<input type="submit" value="submit">
</form>
{{ end }}`