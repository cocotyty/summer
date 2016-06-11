package gen

import (
	"go/ast"
	"bytes"
	"text/template"
	"strings"
)

type SortPlugin struct {
	tpl *template.Template
}

func NewSortPlugin() *SortPlugin {
	return &SortPlugin{template.Must(template.New("name").Parse(sortTpl))}
}

func (SortPlugin *SortPlugin)Tag() string {
	return "sort"
}
func (SortPlugin *SortPlugin)Handle(body *bytes.Buffer, imports *bytes.Buffer, spec *ast.TypeSpec, tag string, path string, pkg string) {
	if _, ok := spec.Type.(*ast.ArrayType); ok {
		sorts := []string{}
		if tag != "*" {
			tags := strings.Split(tag, ",")
			for _, v := range tags {
				sorts = append(sorts, strings.TrimSpace(v))
			}
		}
		SortPlugin.tpl.Execute(body, map[string]interface{}{
			"name":spec.Name,
			"sorts":sorts,
		})
	}
}

var sortTpl = `
func (arr {{.name}})Len()int{
	return len(arr)
}
func (arr {{.name}})Less(i,j int)bool{
	{{range .sorts}}
	if arr[i].{{.}}!=arr[j].{{.}}{
		return arr[i].{{.}}!<arr[j].{{.}}
	}
	{{else}}
	return arr[i]<arr[j]
	{{end}}
	{{with .sorts}}
	return false
	{{end}}
}

func (arr {{.name}})Swap(i, j int){
	arr[i],arr[j]=arr[j],arr[i]
}
`