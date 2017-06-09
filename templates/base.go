package templates

const baseSrc = `<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{ .Title }}</title>
	{{ block "head" . }}{{ end }}
</head>

<body>
	{{ block "content" . }}{{ end }}
</body>
</html>`