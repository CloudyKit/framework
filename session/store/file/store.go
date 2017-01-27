// MIT License
//
// Copyright (c) 2017 José Santos <henrique_1609@me.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package file

import (
	"github.com/CloudyKit/framework/container"
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

func (store store) Reader(_ *container.IoC, name string, after time.Time) (reader io.ReadCloser, err error) {
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

func (store store) Writer(_ *container.IoC, name string) (writer io.WriteCloser, err error) {
	writer, err = os.Create(path.Join(store.BaseDir, name))
	return
}

func (store store) Remove(_ *container.IoC, name string) error {
	return os.Remove(path.Join(store.BaseDir, name))
}

func (store store) GC(_ *container.IoC, before time.Time) {
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
