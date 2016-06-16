package main

import (
	"flag"
	"fmt"
	"github.com/CloudyKit/jet"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	output   = flag.String("o", "cditypes.go", "-o=output.go")
	pkg      = flag.String("p", "", "-p=pkg")
	filename = flag.String("f", "cdi.txt", "-f=cdi.txt")
)

type Type struct {
	Name, TypeName, GetterName string
	IsPtr                      bool
}

func (typ *Type) GetTypeValue() string {
	return fmt.Sprintf("(*%s)(nil)", typ.GetType())
}

func (typ *Type) GetType() string {
	if typ.IsPtr {
		return "*" + typ.Name
	}
	return typ.Name
}

type File struct {
	Pkg       string
	Types     []*Type
	generated bool
	lines     []string
}

var skip = len("//cdi:")

func (types *File) ParseComments(filecontent string) {
	for _, line := range strings.Split(filecontent, "\n") {

		if strings.HasPrefix(line, "///cdi:generated") {
			types.generated = !types.generated
			continue
		}

		if types.generated {
			continue
		}

		types.lines = append(types.lines, line)
		if !strings.HasPrefix(line, "//cdi:") {
			continue
		}

		line := strings.Split(strings.TrimSpace(line[skip:]), " ")
		name := strings.TrimSpace(line[0])

		if name != "" {

			isPtr := name[0] == '*'
			if isPtr {
				name = name[1:]
			}

			var typeName string
			var getterName string

			if len(line) < 3 {
				getterName = "Get" + name
			} else {
				getterName = strings.TrimSpace(line[2])
			}
			if len(line) < 2 {
				typeName = name + "Type"
			} else {
				typeName = strings.TrimSpace(line[1])
			}

			types.Types = append(types.Types, &Type{
				IsPtr:      isPtr,
				Name:       name,
				TypeName:   typeName,
				GetterName: getterName,
			})
		}
	}
}

func main() {
	flag.Parse()
	types := &File{Pkg: *pkg}

	if types.Pkg == "" {
		dir, _ := os.Getwd()
		types.Pkg = path.Base(dir)
	}

	filecontent, err := ioutil.ReadFile(*filename)
	if err != nil {
		panic(err)
		return
	}

	file, err := os.Create(*output)
	if err != nil {
		panic(err)
		return
	}
	defer file.Close()

	template, _ := jet.NewSet().LoadTemplate("--", `{{if .Pkg}}
package {{.Pkg}}
import "github.com/CloudyKit/framework/cdi"
{{end}}
///cdi:generated
{{range .Types}}var {{.TypeName}} = cdi.TypeOfElem({{.GetTypeValue()}})
func {{.GetterName}}(c *cdi.DI) {{.GetType()}} {
	v,_:=c.Val4Type({{.TypeName}}).({{.GetType()}})
	return v
}
{{end}}
///cdi:generated`)

	fset := token.NewFileSet()
	gof, _ := parser.ParseFile(fset, *filename, filecontent, parser.ImportsOnly)

	filecontentstring := string(filecontent)

	if gof == nil {
		types.ParseComments(filecontentstring)
		for _, line := range types.lines {
			file.WriteString(line + "\n")
		}
	} else {
		types.ParseComments(filecontentstring[0:gof.Package])
		types.lines = nil
		file.Write(filecontent[0:gof.End()])

		types.ParseComments(filecontentstring[gof.End():])
		for _, _import := range gof.Imports {
			path, _ := strconv.Unquote(_import.Path.Value)
			if path == "github.com/CloudyKit/framework/cdi" {
				goto found
			}
		}
		file.WriteString(`import "github.com/CloudyKit/framework/cdi"` + "\n")
	found:
		types.Pkg = ""
		for _, line := range types.lines {
			file.WriteString(line + "\n")
		}
	}

	template.Execute(file, nil, types)
}
