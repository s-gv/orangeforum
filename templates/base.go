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
		padding-bottom: 60px;
	}
	#footer {
		position: absolute;
		bottom: 0;
		width: 100%;
		height: 40px;
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
	a.blacklink, a.blacklink:link, a.blacklink:hover, a.blacklink:active, a.blacklink:visited {
		color: #000;
	}
	a.footerlink, a.footerlink:link, a.footerlink:hover, a.footerlink:active, a.footerlink:visited {
		color: grey;
	}
	#navleft {
		float: left;
		max-width: 70%;
	}
	#navright {
		float: right;
	}
	@media screen and (min-width:600px) {
		#forum {
			float: left;
		}
		#nav {
			float: left;
			margin-left: 20px;
		}
	}
	.muted {
		color: grey;
	}
	.row {
		margin-top: 20px;
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
	</style>
	<title>{{ .Common.ForumName }}</title>
	{{ block "head" . }}{{ end }}
</head>

<body>
	<div id="container">
		<div id="header" class="clearfix">
			<div id="navleft">
				<div id="forum"><a href="/" class="blacklink">{{ .Common.ForumName }}</a></div>
				<div id="nav">
				{{ block "navbar" . }}
				{{ end }}
				</div>
			</div>
			<div id="navright">
				{{ if .Common.UserName }}
				<a class="blacklink" href="/users?u={{ .Common.UserName }}">{{ .Common.UserName }}</a> | <a class="blacklink" href="/logout">Logout</a>
				{{ else }}
				<a class="blacklink" href="/login?next={{ .Common.CurrentURL }}">Login</a>
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
		<a class="footerlink" href="/note?id={{ $e.ID }}">{{ $e.Name }}</a>
		{{ end }}
		</div>
	</div>
	{{ .Common.BodyAppendage }}
</body>
</html>`