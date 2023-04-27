package encore

import (
	"encoding/json"
	"github.com/CloudyKit/framework/app"
	"github.com/CloudyKit/framework/request"
	"github.com/CloudyKit/framework/view"
	"os"
	"path"
)

type Component struct {
	BuildPath string
}

type Manifest struct {
	EntryPoints map[string]map[string][]string `json:"entrypoints"`
}

func (component Component) Bootstrap(a *app.Kernel) {
	manifest, err := ReadManifest(component.BuildPath)
	if err != nil {
		panic(err)
	}
	view.GetJetSet(a.Registry).AddGlobal("encoreManifest", manifest)

	a.BindFilterFuncHandlers(func(c *request.Context) {
		manifest, err := ReadManifest(component.BuildPath)
		if err != nil {
			panic(err)
		}
		view.GetJetSet(a.Registry).AddGlobal("encoreManifest", manifest)
		c.Next()
	})
}

func ReadManifest(buildPath string) (manifest *Manifest, err error) {
	var f *os.File
	f, err = os.Open(path.Join(buildPath, "entrypoints.json"))
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&manifest)
	return
}
