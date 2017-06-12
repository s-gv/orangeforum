package templates

const extranoteSrc = `
{{ define "content" }}


<h1>{{ .Name }}</h1>
<p>{{ .UpdatedDate.Format "02 Jan 2006" }}</p>
<div>{{ .Content }}</div>

{{ end }}`