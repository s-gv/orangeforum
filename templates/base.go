package templates

const basesrc = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
	{{ block "head" . }}{{ end }}
</head>

<body>
	{{ block "content" . }}{{ end }}
</body>
</html>`