package templates

const baseSrc = `<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>
	html, body {
		margin: 0;
		padding: 0;
		height: 100%;
	}
	#container {
		max-width: 960px;
		line-height: 1.58;
		margin: 0 auto;
		min-height: 100%;
		position: relative;
		padding-left: 10px;
		padding-right: 10px;
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
	#footer a, #footer a:link, #footer a:hover, #footer a:active, #footer a:visited {
		color: grey;
	}
	.link-btn, .link-btn:link, .link-btn:visited {
		color: white;
		background: #07C;
		padding: 10px 20px;
		text-align: center;
	}
	.link-btn:hover {
		background: #3af;
	}
	#navleft {
		float: left;
		max-width: 70%;
	}
	#navright {
		float: right;
	}
	.muted {
		color: grey;
	}
	.row {
		margin-top: 20px;
	}
	table {
		width: 75%;
	}
	th {
		text-align: right;
	}
	td, th {
		padding: 5px;
	}
	@media (max-width: 599px) {
		table {
			width: 100%;
		}
		th {
			float: none;
			display: block;
			text-align: left;
		}
		td {
			display: block;
			float: none;
		}
	}
	input[type="text"], input[type="number"], textarea {
		width: 100%;
	}
	</style>
	<title>{{ .Common.ForumName }}</title>
	{{ block "head" . }}{{ end }}
</head>

<body>
	<div id="container">
		<div id="header" class="clearfix">
			<div id="navleft">
				<a href="/">{{ .Common.ForumName }}</a>
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