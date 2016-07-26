package file

import (
	"github.com/CloudyKit/framework/scope"
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

func (store store) Reader(_ *scope.Variables, name string, after time.Time) (reader io.ReadCloser, err error) {
	var stat os.FileInfo
	sessionFile := path.Join(store.BaseDir, name)
	stat, err = os.Stat(sessionFile)

	if err == nil {
		if stat.ModTime().After(after) {
			reader, err = os.Open(sessionFile)
		} else {
			os.Remove(sessionFile)
		}
		if err == nil {
			return
		}
	}

	if os.IsNotExist(err) {
		err = nil
	}
	return
}

func (store store) Writer(_ *scope.Variables, name string) (writer io.WriteCloser, err error) {
	writer, err = os.Create(path.Join(store.BaseDir, name))
	return
}

func (store store) Remove(_ *scope.Variables, name string) error {
	return os.Remove(path.Join(store.BaseDir, name))
}

func (store store) GC(_ *scope.Variables, before time.Time) {
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
