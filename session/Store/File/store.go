package file

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type Store struct {
	BaseDir string
}

func (store Store) Reader(name string) (reader io.ReadCloser) {
	var err error
	reader, err = os.Open(path.Join(store.BaseDir, name))
	if err != nil && os.IsNotExist(err) {
		panic(err)
	}
	return
}

func (store Store) Writer(name string) (writer io.WriteCloser) {
	var err error
	writer, err = os.Create(path.Join(store.BaseDir, name))
	if err != nil {
		panic(err)
	}
	return
}

func (store Store) Touch(name string) error {
	now := time.Now()
	return os.Chtimes(name, now, now)
}

func (store Store) Gc(before time.Time) {
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
