// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const pmSrc = `
{{ define "content" }}

<h1 id="title"><a href="/pm">Private Messages</a></h1>

{{ if .Messages }}
{{ range .Messages }}
<div class="comment-row">
	<div class="comment-title muted">
		{{ if not .IsRead }}<span class="alert">&#x2757;</span>{{ end }}
		<a href="/users?u={{ .From }}">{{ .From }}</a> {{ .CreatedDate }} |
		<a href="/pm?quote={{ .ID }}#end">reply</a> |
		<form method="post" action="/pm/delete" style="display: inline;">
			<input type="hidden" name="csrf" value="{{ $.Common.CSRF }}">
  			<input type="hidden" name="id" value="{{ .ID }}">
			<input type="hidden" name="lmd" value="{{ $.FirstMessageDate }}">
  			<button type="submit" name="submit" value="delete" class="link-button">delete</button>
		</form>
	</div>
	<div class="comment">{{ .Content }}</div>
</div>
<hr class="sep">
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No messages to show.</div>
</div>
{{ end }}
{{ if .LastMessageDate }}
<a href="/pm?lmd={{ .LastMessageDate }}">More</a>
{{ end }}

<h2 id="end" style="margin-top: 40px;">Send Message</h2>
<div>
<form action="/pm/new" method="POST">
	<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
	<table class="form">
		<tr>
			<th>To:</th>
			<td><input type="text" name="to" placeholder="username" value="{{ .To }}"></td>
		</tr>
		<tr>
			<th>Content:</th>
			<td><textarea name="content" rows="12">{{ .Content }}</textarea></td>
		</tr>
		{{ if .Common.Msg }}
		<tr>
			<th></th>
			<td><span class="alert">{{ .Common.Msg }}</span></td>
		</tr>
		{{ end }}
		<tr>
			<th></th>
			<td><input type="submit" name="action" class="no-double-post" value="Send message"></td>
		</tr>
	</table>
</form>
</div>

{{ end }}
`
