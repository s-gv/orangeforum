// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates


const commenteditSrc = `
{{ define "content" }}

<h2 id="title"><a href="/topics?id={{ .TopicID }}">{{ .TopicName }}</a></h2>
<p id="subtitle" class="muted"><a href="/users?u={{ .TopicOwnerName }}">{{ .TopicOwnerName }}</a> in <a href="/groups?name={{ .GroupName }}">{{ .GroupName }}</a> {{ .TopicCreatedDate }}</p>

<div>{{ .ParentComment }}</div>

<div class="row">
	<form action="{{ if .CommentID }}/comments/edit{{ else }}/comments/new{{ end }}" method="POST" enctype="multipart/form-data">
	<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
	<input type="hidden" name="id" value="{{ .CommentID }}">
	<input type="hidden" name="tid" value="{{ .TopicID }}">
	<textarea name="content" rows="12">{{ .Content }}</textarea>

	{{ if .IsImageUploadEnabled }}
	<div>Add Image (optional): <input type="file" name="img" accept="image/*"></div>
	{{ end }}

	{{ if or .IsMod .IsAdmin .IsSuperAdmin }}
	<div><input type="checkbox" name="is_sticky"{{ if .IsSticky }} checked{{ end }}> Sticky</div>
	{{ end }}

	{{ .Common.Msg }}

	<div>
	{{ if .CommentID }}
		{{ if not .IsDeleted }}
		<input type="submit" name="action" value="Update">
		<input type="submit" name="action" value="Delete">
		{{ else }}
		<input type="submit" name="action" value="Undelete">
		{{ end }}
	{{ else }}
	<input type="submit" name="action" value="Submit reply">
	{{ end }}
	</div>

	</form>
</div>

{{ end }}`
