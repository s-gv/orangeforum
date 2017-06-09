package templates

const adminindexSrc = `
{{ define "content" }}

Yo
{{ index .Config "ForumName" }}

<form action="/admin" method="POST">

<input type="hidden" name="csrf" value="{{ .CSRF }}">
Forum Name: <input type="text" name="forum_name" value="{{ index .Config "forum_name" }}" required>
Announcement: <input type="text" name="header_msg" value="{{ index .Config "header_msg" }}">
SMTP Host: <input type="text" name="smtp_host" value="{{ index .Config "smtp_host" }}">
SMTP Port: <input type="number" name="smtp_port" value="{{ index .Config "smtp_port" }}">
SMTP Username: <input type="text" name="smtp_user" value="{{ index .Config "smtp_user" }}">
SMTP Password: <input type="text" name="smtp_pass" value="{{ index .Config "smtp_pass" }}">

<input type="submit" value="Update">

</form>

{{ end }}`