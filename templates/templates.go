package templates

import (
	"html/template"
	"io"
)

var tmpls map[string]*template.Template = make(map[string]*template.Template)

func init() {
	tmpls["index.html"] = template.Must(template.New("base").Parse(basesrc))
	template.Must(tmpls["index.html"].New("index").Parse(indexsrc))

	tmpls["signup.html"] = template.Must(template.New("base").Parse(basesrc))
	template.Must(tmpls["signup.html"].New("signup").Parse(signupsrc))

	tmpls["login.html"] = template.Must(template.New("base").Parse(basesrc))
	template.Must(tmpls["login.html"].New("login").Parse(loginsrc))
}

func Render(wr io.Writer, template string, data interface{}) error {
	return tmpls[template].Execute(wr, data)
}