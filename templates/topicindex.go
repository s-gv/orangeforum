// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const topicindexSrc = `
{{ define "content" }}

<div class="btn-row">
	{{ if not .IsClosed }}
	<a class="link-btn" href="/comments/new?tid={{ .TopicID }}">Reply</a>
	{{ end }}
	{{ if or .IsAdmin .IsMod .IsSuperAdmin (and .IsOwner (not .IsClosed)) }}
	<a class="link-btn" href="/topics/edit?id={{ .TopicID }}">Edit topic</a>
	{{ end }}
	{{ if and .Common.UserName .Common.IsTopicSubAllowed }}
	{{ if .SubToken }}
	<form action="/topics/unsubscribe?token={{ .SubToken }}" method="POST">
		<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
		<input class="btn" type="submit" value="Unsubscribe">
	</form>
	{{ else }}
	<form action="/topics/subscribe?id={{ .TopicID }}" method="POST">
		<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
		<input class="btn" type="submit" value="Subscribe">
	</form>
	{{ end }}
	{{ end }}
</div>

<h2 id="title"><a href="/topics?id={{ .TopicID }}">{{ .TopicName }}{{ if .IsClosed }} [closed]{{ end }}</a></h2>
<div class="comment-title muted"><a href="/users?u={{ .OwnerName }}">{{ .OwnerName }}</a> in <a href="/groups?name={{ .GroupName }}">{{ .GroupName }}</a> {{ .CreatedDate }}</div>
<div class="comment-row">
	<div class="comment">
		<p>{{ .Content }}</p>
	</div>
</div>
<hr class="sep">

{{ if .Comments }}
{{ range .Comments }}
<div class="comment-row" id="comment-{{ .ID }}">
	<div class="comment-title muted"><a href="/users?u={{ .UserName }}">{{ .UserName }}</a> <a href="/comments?id={{ .ID }}">{{ .CreatedDate }}</a>{{ if or .IsOwner $.IsAdmin $.IsMod $.IsSuperAdmin }} | <a href="/comments/edit?id={{ .ID }}">edit</a>{{end}} | <a href="/comments/new?tid={{ $.TopicID }}&quote={{ .ID }}">quote</a></div>
	{{ if .IsDeleted }}
		<div class="comment">[DELETED]</div>
	{{ else }}
		<div class="comment">{{ .Content }}</div>
		{{ if .ImgSrc }}<div><img src="/img?name={{ .ImgSrc }}"></div>{{ end }}
	{{ end }}
</div>
<hr class="sep">
{{ end }}
{{ else }}
<div class="row">
	<div class="muted">No comments to show.</div>
</div>
{{ end }}
<div id="comment-last"></div>

{{ if gt .NumPages 1 }}
	<div style="float: right; max-width: 70%;">
	Pages:
	{{ range $i, $e := .Pages }}
		{{ if eq $i $.CurrentPage }}
		{{ $i }}
		{{ else }}
		<a href="/topics?id={{ $.TopicID }}&p={{ $i }}">{{ $i }}</a>
		{{ end }}
	{{ end }}
	</div>
{{ end }}

{{ if not .IsLastPage }}
<div>
	<div><a href="/topics?id={{ .TopicID }}&p={{ .NextPage }}">Next Page</a></div>
</div>
{{ end }}

{{ if .Common.UserName }}
<div style="margin-top: 40px;">
<form action="/comments/new" method="POST" enctype="multipart/form-data">
	<input type="hidden" name="csrf" value="{{ .Common.CSRF }}">
	<input type="hidden" name="id" value="{{ .CommentID }}">
	<input type="hidden" name="tid" value="{{ .TopicID }}">
	<textarea name="content" rows="12" placeholder="Your comment..."></textarea>
	{{ if .IsImageUploadEnabled }}
	<div>Add Image (optional): <input type="file" name="img" accept="image/*"></div>
	{{ end }}
	<input type="submit" name="action" class="no-double-post" value="Add comment">
</form>
</div>
{{ end }}

{{ end }}`
