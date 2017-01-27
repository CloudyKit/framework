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

package container

import (
	"io"
	"sync"
)

var (
	//pools
	ioCloserProviderPool = sync.Pool{
		New: func() interface{} {
			return new(ioCloserProvider)
		},
	}

	poolerProviderPool = sync.Pool{
		New: func() interface{} {
			return new(poolerProvider)
		},
	}
)

func NewIOCloserProvider(v io.Closer) (closer *ioCloserProvider) {
	closer, _ = ioCloserProviderPool.Get().(*ioCloserProvider)
	closer.Value = v
	return
}

func NewPoolProvider(pool *sync.Pool, v interface{}) (pooler *poolerProvider) {
	pooler, _ = poolerProviderPool.Get().(*poolerProvider)
	pooler.Pool = pool
	pooler.Value = v
	return
}

type ioCloserProvider struct {
	Value io.Closer
}

type poolerProvider struct {
	Pool  *sync.Pool
	Value interface{}
}

func (pooler *poolerProvider) Provide(c *IoC) interface{} {
	if pooler.Value != nil {
		return pooler.Value
	}
	pooler.Value = pooler.Pool.Get()
	return pooler.Value
}

func (pooler *poolerProvider) Dispose() {
	if pooler.Value != nil {
		pooler.Pool.Put(pooler.Value)
	}
	poolerProviderPool.Put(pooler)
}

func (pp *ioCloserProvider) Dispose() {
	closer := pp.Value
	ioCloserProviderPool.Put(pp)
	closer.Close()
}

func (pp *ioCloserProvider) Provide(_ *IoC) interface{} {
	return pp.Value
}
