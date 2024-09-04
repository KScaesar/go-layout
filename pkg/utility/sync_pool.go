package utility

import (
	"sync"
)

func NewPool[T any](newFn func() *T) *Pool[T] {
	return &Pool[T]{
		syncPool: sync.Pool{
			New: func() any { return newFn() },
		},
	}
}

type Pool[T any] struct {
	syncPool sync.Pool
}

func (p *Pool[T]) Get() *T {
	return p.syncPool.Get().(*T)
}

func (p *Pool[T]) Put(value *T) {
	p.syncPool.Put(value)
}
