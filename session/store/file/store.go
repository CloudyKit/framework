package file

import (
	"github.com/CloudyKit/framework/context"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

type store struct {
	BaseDir string
}

func New(directory string) store {
	directory, _ = filepath.Abs(directory)
	_, err := os.Stat(directory)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(directory, 0666)
	}
	return store{directory}
}

func (store store) Reader(_ *context.Context, name string) (reader io.ReadCloser, err error) {
	reader, err = os.Open(path.Join(store.BaseDir, name))
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return
}

func (store store) Writer(_ *context.Context, name string) (writer io.WriteCloser, err error) {
	writer, err = os.Create(path.Join(store.BaseDir, name))
	return
}

func (store store) Remove(_ *context.Context, name string) error {
	return os.Remove(path.Join(store.BaseDir, name))
}

func (store store) Gc(_ *context.Context, before time.Time) {
	files, err := ioutil.ReadDir(store.BaseDir)
	if err != nil {
		panic(err)
	}
	numFiles := len(files)
	for i := 0; i < numFiles; i++ {
		file := files[i]
		if !file.IsDir() && file.ModTime().Before(before) {
			os.Remove(path.Join(store.BaseDir, file.Name()))
		}
	}
}
