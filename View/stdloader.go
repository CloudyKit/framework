package View

import (
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"text/template/parse"
)

func NewStdTemplateLoader(base string) *StdTemplateLoader {
	stdLoader := new(StdTemplateLoader)
	stdLoader.BaseDir = base
	stdLoader.baseTemplate = template.New("baseStdTemplate")
	stdLoader.Funcs = make(template.FuncMap)
	return stdLoader
}

type StdTemplateLoader struct {
	BaseDir      string
	Funcs        template.FuncMap
	baseTemplate *template.Template
}

func (stdLoader *StdTemplateLoader) Refresh() {
	stdLoader.baseTemplate = template.New("baseStdTemplate")
}

func (stdLoader *StdTemplateLoader) autoLoad(list *parse.ListNode) {
	if list != nil {
		for i := 0; i < len(list.Nodes); i++ {
			switch node := list.Nodes[i].(type) {
			case *parse.TemplateNode:
				stdLoader.View(node.Name)
			case *parse.BranchNode:
				stdLoader.autoLoad(node.List)
				stdLoader.autoLoad(node.ElseList)
			case *parse.IfNode:
				stdLoader.autoLoad(node.List)
				stdLoader.autoLoad(node.ElseList)
			case *parse.RangeNode:
				stdLoader.autoLoad(node.List)
				stdLoader.autoLoad(node.ElseList)
			case *parse.WithNode:
				stdLoader.autoLoad(node.List)
				stdLoader.autoLoad(node.ElseList)
			}
		}
	}
}

type stdRender template.Template

func (tt *stdRender) Execute(w io.Writer, c Data) error {
	return (*template.Template)(tt).Execute(w, c)
}

func (stdLoader *StdTemplateLoader) View(name string) (view ViewRenderer, err error) {
	t := stdLoader.baseTemplate.Lookup(name)
	if t == nil {
		var b []byte
		b, err = ioutil.ReadFile(filepath.Join(stdLoader.BaseDir, name))
		if err != nil {
			return
		}
		t, err = stdLoader.baseTemplate.New(name).Funcs(stdLoader.Funcs).Parse(string(b))
		if err == nil {
			stdLoader.autoLoad(t.Tree.Root)
		}
	}
	view = (*stdRender)(t)
	return
}
