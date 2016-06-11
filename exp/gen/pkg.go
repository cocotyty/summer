package gen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"reflect"
	"bytes"
	"github.com/cocotyty/summer"
)
var log = summer.NewSimpleLog("gen",summer.InfoLevel)
var L = string(os.PathSeparator)

type Plugin interface {
	Tag() string
	Handle(writer *bytes.Buffer, imports *bytes.Buffer, spec *ast.TypeSpec, tag string, path string, pkg string)
}
type PluginList []Plugin
type GenCode struct {
	RootPath   string
	WorkPath   string
	PluginList PluginList
	pkgs       map[string]*ast.Package
}

func (code *GenCode) Init(workPath string) {
	path := os.Getenv("GOPATH")
	log.Println(path)
	if os.PathSeparator == '\\' {
		path = strings.Replace(path, "/", "\\", -1)
	}
	gopaths := strings.Split(path, string(os.PathListSeparator))
	code.WorkPath = workPath
	log.Println(workPath)
	for _, v := range gopaths {
		if strings.Contains(workPath, v) {
			code.RootPath = v + L + "src" + L
		}
	}
	if code.RootPath == "" {
		panic("what fuck?")
	}
	code.pkgs = map[string]*ast.Package{}
}
func (code *GenCode) Work(path string) {
	log.Println("work")
	fs := token.NewFileSet()
	pkgs, _ := parser.ParseDir(fs, path, func(f os.FileInfo) bool {
		return true
	}, parser.ParseComments)
	for _, pkg := range pkgs {
		body := bytes.NewBuffer(nil)
		header := bytes.NewBuffer(nil)
		header.WriteString("package " + pkg.Name + "\n")
		header.WriteString("import (\n")
		for _, f := range pkg.Files {
			log.Println("f:", f)
			for _, d := range f.Decls {
				log.Println("d:", d)
				switch dd := d.(type){
				case *ast.GenDecl:
					log.Println("dd:", dd)
					for _, spec := range dd.Specs {
						switch sp := spec.(type) {
						case *ast.TypeSpec :
							log.Println(sp.Doc, sp.Comment, sp.Type)
							log.Println(sp.Name, dd.Doc)
							if dd.Doc != nil {
								for _, doc := range dd.Doc.List {
									text := strings.TrimSpace(doc.Text)
									if strings.HasPrefix(text, "//") {
										text = strings.TrimSpace(text[2:])
									}
									if strings.HasPrefix(text, "SM ") {
										code.PluginWork(body, header, sp, text[3:], path, pkg.Name)
										break
									}
								}
							}
						}
					}
				}
			}
		}
		header.WriteString("\n)")
		f, err := os.Create(path + L + "go_gen.go")
		if err != nil {
			log.Println(err)
			return
		}
		header.Write(body.Bytes())
		f.Write(header.Bytes())
		f.Close()
	}
}
func (code *GenCode) Watch() {

}

func (code *GenCode)PluginWork(body, header *bytes.Buffer, spec *ast.TypeSpec, smTag string, path string, pkg string) {
	tags := reflect.StructTag(smTag)

	for _, v := range code.PluginList {
		log.Println(v.Tag(), tags)
		tag := tags.Get(v.Tag())
		if tag != "" {
			v.Handle(body, header, spec, tag, path, pkg)
		}
	}

	//switch  t := spec.Type.(type) {
	//case *ast.StructType:
	//	log.Println("struct", t)
	//case *ast.ArrayType:
	//	log.Println("array")
	//case *ast.ChanType:
	//	log.Println("chan")
	//case *ast.InterfaceType:
	//	log.Println("interface")
	//case *ast.FuncType:
	//	log.Println("func")
	//case *ast.MapType:
	//	log.Println("map")
	//}
}
func (code *GenCode)Register(Plugin Plugin) {
	code.PluginList = append(code.PluginList, Plugin)
}
// func Pkg(name string) PkgHandler {
//
// 	log.Println(gopaths)
// 	fs := token.NewFileSet()
// 	var pkgs map[string]*ast.Package
// 	var err error
// 	for _, gopath := range gopaths {
// 		pkgs, err = parser.ParseDir(fs, gopath+string(os.PathSeparator)+"src"+string(os.PathSeparator)+name, func(f os.FileInfo) bool {
// 			return true
// 		}, parser.ParseComments)
// 		log.Println(gopath + string(os.PathSeparator) + "src" + string(os.PathSeparator) + name)
// 		if err == nil {
// 			break
// 		}
// 	}
// 	if err != nil {
// 		log.Println(err)
// 		return nil
// 	}
// 	for _, pkg := range pkgs {
// 		return &pkgHandler{name, pkg}
// 	}
// 	return nil
// }

// type pkgHandler struct {
// 	name string
// 	pkg  *ast.Package
// }

// func (this *pkgHandler) Name() string {
// 	return this.name
// }
// func (this *pkgHandler) StructList() map[string]StructHandler {
// 	for _, file := range this.pkg.Files {
// 		for _, decl := range file.Decls {
// 			if k, ok := decl.(*ast.FuncDecl); ok {
// 				log.Println("func:", k.Name, k.Doc)
// 				if k.Type.Params != nil {
// 					log.Println("params:")
// 					for _, filed := range k.Type.Params.List {
// 						log.Println(filed.Names, " ", filed.Type)
// 					}
// 				}
// 				if k.Type.Results != nil {
// 					log.Println("results:")
// 					for _, filed := range k.Type.Results.List {
// 						log.Println(filed.Names, " ", filed.Type)
// 					}
// 				}
// 			}
// 			if k, ok := decl.(*ast.GenDecl); ok {
// 				for _, spec := range k.Specs {
// 					if sp, ok := spec.(*ast.ValueSpec); ok {
// 						for _, s := range sp.Names {
// 							log.Println("var:", s.Name)
// 						}
// 					}
// 					if sp, ok := spec.(*ast.TypeSpec); ok {
// 						log.Println("type:", sp.Name)
// 					}
// 					if sp, ok := spec.(*ast.ImportSpec); ok {
// 						log.Println("import:", sp.Name, sp.Path)
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }
// func (this *pkgHandler) InterfaceList() map[string]InterfaceHandler { return nil }
// func (this *pkgHandler) Func(name string) FuncHandler               { return nil }
// func (this *pkgHandler) FuncList() map[string]FuncHandler           { return nil }
