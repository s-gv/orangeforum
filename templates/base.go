// Copyright (c) 2017 Sagar Gubbi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package templates

const baseSrc = `<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
	* {
		-webkit-box-sizing: border-box;
		-moz-box-sizing: border-box;
		box-sizing: border-box;
	}
	html, body {
		margin: 0;
		height: 100%;
	}
	#container {
		max-width: 960px;
		line-height: 1.58;
		margin: 0 auto;
		min-height: 100%;
		position: relative;
	}
	#header {
		padding-top: 10px;
	}
	#content {
		clear: both;
		padding-top: 20px;
		padding-bottom: 75px;
	}
	#footer {
		position: absolute;
		bottom: 0;
		width: 100%;
		height: 60px;
		text-align: center;
	}
	.clearfix {
		overflow: auto;
	}
	body {
		font-family: Arial, "Helvetica Neue", Helvetica, sans-serif;
		text-rendering: optimizeLegibility;
		-webkit-font-smoothing: antialiased;
		padding-left: 10px;
		padding-right: 10px;
	}
	img {
		max-width: 100%;
		display: block;
		margin: 0 auto;
	}
	a {
		text-decoration: none;
	}
	a:link {
		color: #07C;
	}
	a:hover, a:active {
		color: #3af;
	}
	a:visited {
		color: #005999;
	}
	#header a, #header a:link, #header a:hover, #header a:active, #header a:visited {
		color: #000;
	}
	#title a, #title a:link, #title a:hover, #title a:active, #title a:visited {
		color: #000;
	}
	#footer a, #footer a:link, #footer a:hover, #footer a:active, #footer a:visited {
		color: grey;
	}
	.link-btn, .link-btn:link, .link-btn:visited {
		color: white;
		background: #07C;
		padding: 10px 20px;
		text-align: center;
		width: 150px;
		margin-left: 20px;
		font-size: 16px;
	}
	.link-btn:hover {
		background: #3af;
	}
	.btn {
		padding: 10px 20px;
		background: #07C;
		font-size: 16px;
		color: white;
		border: none;
		margin-left: 20px;
		width: 150px;
		line-height: inherit;
	}
	.btn:hover {
		background: #3af;
		cursor: pointer;
		cursor: hand;
	}
	.btn-row form, .btn-row a {
		display: inline-block;
	}
	.btn-row {
		text-align: right;
	}
	@media (max-width: 599px) {
		.btn {
			font-size: 12px;
			padding: 5px 10px;
			margin: 10px;
			width: 100px;
		}
		.link-btn, .link-btn:link, .link-btn:visited {
			font-size: 12px;
			padding: 5px 10px;
			margin: 10px;
			width: 100px;
		}
		.btn-row {
			text-align: center;
		}
	}
	#navleft {
		float: left;
		max-width: 70%;
	}
	#navright {
		float: right;
	}
	.muted {
		color: darkgrey;
	}
	.muted a, .muted a:link, .muted a:hover, .muted a:visited, .muted a:active {
		color: grey;
	}
	.row {
		margin-top: 20px;
	}
	th, td {
		text-align: left;
	}
	@media (min-width: 600px) {
		.form td {
			width: 360px;
		}
	}
	.form th {
		text-align: right;
		padding-top: 10px;
	}
	.form td {
		padding-top: 10px;
	}
	@media (max-width: 599px) {
		table.form {
			width: 100%;
		}
		.form th {
			display: block;
			float: left;
			text-align: left;
			max-width: 75%;
			padding-top: 20px;
		}
		.form td {
			display: block;
			text-align: right;
			padding-top: 20px;
		}
		.form input[type="submit"] {
			width: 100%;
		}
	}
	.form input[type="text"], .form input[type="number"], .form input[type="email"], .form input[type="password"] {
		width: 100%;
	}
	textarea {
		width: 100%;
	}
	.sep {
		border: none;
		height: 1px;
		background-color: #ccc;
	}
	#title {
		margin-bottom: 0px;
	}
	#subtitle {
		margin-top: 0px;
	}
	.comment p:first-child {
		margin-top: 0px;
	}
	</style>
	<title>{{ .Common.ForumName }}</title>
	{{ block "head" . }}{{ end }}
</head>

<body>
	<div id="container">
		<div id="header" class="clearfix">
			<div id="navleft">
				<a href="/">{{ .Common.ForumName }}</a>{{ if .GroupName }} &gt; <a href="/groups?name={{ .GroupName }}">{{ .GroupName }}</a>{{ end }}
			</div>
			<div id="navright">
				{{ if .Common.UserName }}
				<a href="/users?u={{ .Common.UserName }}">{{ .Common.UserName }}</a>
				{{ else }}
				<a href="/login?next={{ .Common.CurrentURL }}">Login</a>
				{{ end }}
			</div>
		</div>
		<hr>
		<div id="content">
		{{ block "content" . }}{{ end }}
		</div>
		<div id="footer">
		{{ range $i, $e := .Common.ExtraNotesShort }}
			{{ if $i }}&middot;{{ end }}
			<a href="/note?id={{ $e.ID }}">{{ $e.Name }}</a>
		{{ end }}
		</div>
	</div>
	{{ .Common.BodyAppendage }}
</body>
</html>`
