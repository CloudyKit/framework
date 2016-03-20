package View

import (
	"errors"
	"github.com/CloudyKit/framework/App"
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Request"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template/parse"
)

var DefaultManager = &Manager{}

var DefaultStdTemplateLoader = NewStdTemplateLoader("./views")

func init() {
	Di.Walkable(Context{})
	App.Default.Put(DefaultManager)

	App.Default.Set((Table)(nil), func(c Di.Context) interface{} {
		tt := tablePool.Get()
		c.Put(tt)
		return tt
	})

	DefaultManager.AddLoader(DefaultStdTemplateLoader, ".go.html", ".html.go")
}

func NewStdTemplateLoader(base string) *StdTemplateLoader {
	stdLoader := new(StdTemplateLoader)
	stdLoader.BaseDir = base
	stdLoader.baseTemplate = template.New("baseStdTemplate")
	stdLoader.Funcs = make(template.FuncMap)
	return stdLoader
}

type Table map[string]interface{}

var tablePool = sync.Pool{
	New: func() interface{} {
		return make(Table)
	},
}

func (t Table) Done() {
	tablePool.Put(t)
}

type Context struct {
	Manager *Manager
	Context *Request.Context
	Data    Table
}

type RendererList struct {
	List []Renderer
}

func (r *RendererList) Append(rs ...Renderer) {
	r.List = append(r.List, rs...)
}

func (r RendererList) Render(c Context) error {
	for i := 0; i < len(r.List); i++ {
		if err := r.List[i].Render(c); err != nil {
			return err
		}
	}
	return nil
}

type Renderer interface {
	Render(Context) error
}

func (r Context) Renderer(v Renderer) error {
	return v.Render(r)
}

func (r Context) Render(view string, context interface{}) error {
	return r.Manager.Render(r.Context.Rw, view, context)
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
	view = t
	return
}

type ViewRenderer interface {
	Execute(w io.Writer, c interface{}) error
}

type viewLoader interface {
	View(name string) (ViewRenderer, error)
}

type viewHandler struct {
	ext string
	viewLoader
}

type viewHandlers []viewHandler

func (s viewHandlers) Len() int {
	return len(s)
}

func (s viewHandlers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s viewHandlers) Less(i, j int) bool {
	return len(s[i].ext) > len(s[j].ext)
}

type Manager struct {
	loaders viewHandlers
}

func (vm *Manager) Render(w io.Writer, name string, context interface{}) (err error) {
	var view ViewRenderer
	view, err = vm.getView(name)
	if err == nil {
		err = view.Execute(w, context)
	}
	return
}

func (vm *Manager) getView(name string) (ViewRenderer, error) {
	for i := 0; i < len(vm.loaders); i++ {
		if strings.HasSuffix(name, vm.loaders[i].ext) {
			return vm.loaders[i].View(name)
		}
	}
	return nil, errors.New("View not found!")
}

func (vm *Manager) AddLoader(loader viewLoader, exts ...string) {
	for i := 0; i < len(exts); i++ {
		vm.loaders = append(vm.loaders, viewHandler{ext: exts[i], viewLoader: loader})
	}
	sort.Sort(vm.loaders)
}
