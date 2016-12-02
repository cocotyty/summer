package summer

import (
	"bytes"
	"text/template"
)

func preTag(obj interface{}, tag string) string {
	t := template.Must(template.New("").Delims("(", ")").Parse(tag))
	buf := bytes.NewBuffer(nil)
	t.Execute(buf, obj)
	return buf.String()
}
