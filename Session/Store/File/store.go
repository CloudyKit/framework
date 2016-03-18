package File

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

func (store Store) Reader(name string) (reader io.ReadCloser, err error) {
	reader, err = os.Open(path.Join(store.BaseDir, name))
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return
}

func (store Store) Writer(name string) (writer io.WriteCloser, err error) {
	writer, err = os.Create(path.Join(store.BaseDir, name))
	return
}

func (store Store) Gc(before time.Time) error {
	files, err := ioutil.ReadDir(store.BaseDir)
	if err != nil {
		return err
	}
	numFiles := len(files)
	for i := 0; i < numFiles; i++ {
		file := files[i]
		if !file.IsDir() && file.ModTime().Before(before) {
			os.Remove(path.Join(store.BaseDir, file.Name()))
		}
	}
	return nil
}
