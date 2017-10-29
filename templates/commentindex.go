// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const commentindexSrc = `
{{ define "content" }}

<div class="row">
	<div class="muted">
		comment by <a href="/users?u={{ .OwnerName }}">{{ .OwnerName }}</a> to <a href="/topics?id={{ .TopicID }}">{{ .TopicName }}</a>
		<a href="/comments?id={{ .ID }}">{{ .CreatedDate }}</a>
		{{ if or .IsOwner $.IsAdmin $.IsMod $.IsSuperAdmin }} | <a href="/comments/edit?id={{ .ID }}">edit</a> {{end}}
	</div>
	{{ if .IsDeleted }}
		<div>[DELETED]</div>
	{{ else }}
		<div class="comment">{{ .Content }}</div>
		{{ if .ImgSrc }}<div><img src="/img?name={{ .ImgSrc }}"></div>{{ end }}
	{{ end }}
</div>


{{ end }}`
