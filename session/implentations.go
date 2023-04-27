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

package session

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"github.com/CloudyKit/framework/container"
	"io"
	"time"
)

// IdGenerator this interface represents an id generator
type IdGenerator interface {
	Generate(id, name string) string
}

type Serializer interface {
	Serialize(session interface{}, w io.Writer) error
	Unserialize(session interface{}, r io.Reader) error
}

type Store interface {
	Reader(c *container.Registry, name string, after time.Time) (io.ReadCloser, error)
	Writer(c *container.Registry, name string) (io.WriteCloser, error)
	Remove(c *container.Registry, name string) error
	GC(c *container.Registry, before time.Time)
}

type RandGenerator struct{}
type GobSerializer struct{}

func (RandGenerator) Generate(id, name string) string {
	if id == "" {
		b := make([]byte, 16)
		n, err := io.ReadFull(rand.Reader, b)
		if n != len(b) || err != nil {
			panic(err)
		}
		return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
	}
	return id
}

func (serializer GobSerializer) Unserialize(dst interface{}, reader io.Reader) error {
	return gob.NewDecoder(reader).Decode(dst)
}

func (serializer GobSerializer) Serialize(src interface{}, writer io.Writer) error {
	return gob.NewEncoder(writer).Encode(src)
}
