package View

import (
	"errors"
	"github.com/CloudyKit/framework/App"
	"github.com/CloudyKit/framework/Di"
	"github.com/CloudyKit/framework/Request"
	"io"
	"sort"
	"strings"
	"sync"
)

var DefaultManager = &Manager{}

var DefaultStdLoader = NewStdTemplateLoader("./views")

func init() {
	Di.Walkable(Context{})
	App.Default.Put(DefaultManager)
	App.Default.Set((Table)(nil), func(c Di.Context) interface{} {
		tt := tablePool.Get()
		c.Put(tt)
		return tt
	})
	DefaultManager.AddLoader(DefaultStdLoader, ".tpl", ".tpl.html")
}

type Table map[string]interface{}

var tablePool = sync.Pool{
	New: func() interface{} {
		return make(Table)
	},
}

func (t Table) Finalize() {
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

func (r Context) Render(view string, context Table) error {
	return r.Manager.Render(r.Context.Rw, view, context)
}

type ViewRenderer interface {
	Execute(w io.Writer, c Table) error
}

type ViewLoader interface {
	View(name string) (ViewRenderer, error)
}

type viewHandler struct {
	ext string
	ViewLoader
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

func (vm *Manager) Render(w io.Writer, name string, context Table) (err error) {
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

func (vm *Manager) AddLoader(loader ViewLoader, exts ...string) {
	for i := 0; i < len(exts); i++ {
		vm.loaders = append(vm.loaders, viewHandler{ext: exts[i], ViewLoader: loader})
	}
	sort.Sort(vm.loaders)
}
