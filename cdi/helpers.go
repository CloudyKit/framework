package cdi

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

func (pooler *poolerProvider) Provide(c *Global) interface{} {
	if pooler.Value != nil {
		return pooler.Value
	}
	pooler.Value = pooler.Pool.Get()
	return pooler.Value
}

func (pooler *poolerProvider) Finalize() {
	if pooler.Value != nil {
		pooler.Pool.Put(pooler.Value)
	}
	poolerProviderPool.Put(pooler)
}

func (pp *ioCloserProvider) Finalize() {
	closer := pp.Value
	ioCloserProviderPool.Put(pp)
	closer.Close()
}

func (pp *ioCloserProvider) Provide(_ *Global) interface{} {
	return pp.Value
}
