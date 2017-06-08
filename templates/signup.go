package templates

const signupSrc = `
{{ define "content" }}

<form action="/signup" method="POST">
<input type="hidden" name="csrf" value="{{ .CSRF }}">
<input type="hidden" name="next" value="{{ .next }}">
Username: <input type="text" name="username" required>
Password: <input type="password" name="passwd" required>
Confirm password: <input type="password" name="confirm" required>
Email (optional): <input type="text" name="email">
Already have an account? <a href="/login?next={{ .next }}">Login</a>
<input type="submit" value="Signup">
</form>

{{ .Msg }}

{{ end }}`