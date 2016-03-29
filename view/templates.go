package view

import (
	"errors"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/common"
	"github.com/CloudyKit/framework/context"
	"github.com/CloudyKit/framework/request"
	"io"
	"reflect"
	"sort"
	"strings"
	"sync"
)

var DefaultManager = &Manager{}

var DefaultStdLoader = NewStdTemplateLoader("./views")

func init() {

	app.Default.Di.Map(DefaultManager)

	app.Default.Di.MapType((*Context)(nil), func(c *context.Context) interface{} {
		tt := DefaultManager.NewContext(c)
		c.Map(tt)
		return tt
	})

	Available(DefaultManager, (*request.Context)(nil))
	AvailableKey(DefaultManager, "linker", (*common.URLer)(nil))
	DefaultManager.AddLoader(DefaultStdLoader, ".tpl", ".tpl.html")
}

func (m *Manager) NewContext(c *context.Context) *Context {
	tt := contextPool.Get().(*Context)

	// init
	tt.Manager = m
	tt.Context = c.Get(tt.Context).(*request.Context)

	// update auto variables
	for i := 0; i < len(m.injectables); i++ {
		tt.Data[m.injectables[i].name] = c.Val4Type(m.injectables[i].typ)
	}

	return tt
}

type Data map[string]interface{}

var contextPool = sync.Pool{
	New: func() interface{} {
		return &Context{Data: make(Data)}
	},
}

type Context struct {
	Manager *Manager //Manager is
	Context *request.Context
	Data    Data
}

func (t *Context) Finalize() {
	contextPool.Put(t)
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

func (r Context) Render(view string, context Data) error {
	return r.Manager.Render(r.Context.Response, view, context)
}

type ViewRenderer interface {
	Execute(w io.Writer, c Data) error
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

type autoset struct {
	name string
	typ  reflect.Type
}

type Manager struct {
	loaders     viewHandlers
	injectables []autoset
}

func (vm *Manager) Render(w io.Writer, name string, context Data) (err error) {
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
