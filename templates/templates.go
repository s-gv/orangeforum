package templates

import (
	"html/template"
	"io"
)

var tmpls map[string]*template.Template = make(map[string]*template.Template)

func init() {
	tmpls["index.html"] = template.Must(template.New("base").Parse(basesrc))
	template.Must(tmpls["index.html"].New("index").Parse(indexsrc))
}

func Render(wr io.Writer, template string, data interface{}) error {
	return tmpls[template].Execute(wr, data)
}