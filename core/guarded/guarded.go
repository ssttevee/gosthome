package guarded

import "sync"

type Value[T any] struct {
	mux sync.Mutex
	val T
}

func New[T any](val T) *Value[T] {
	return &Value[T]{
		mux: sync.Mutex{},
		val: val,
	}
}

func (g *Value[T]) Do(f func(*T)) {
	g.mux.Lock()
	defer g.mux.Unlock()
	f(&g.val)
}

func (g *Value[T]) DoErr(f func(*T) error) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	return f(&g.val)
}

type RWValue[T any] struct {
	mux sync.RWMutex
	val T
}

func NewRW[T any](val T) *RWValue[T] {
	return &RWValue[T]{
		mux: sync.RWMutex{},
		val: val,
	}
}

func (g *RWValue[T]) Write(f func(*T)) {
	g.mux.Lock()
	defer g.mux.Unlock()
	f(&g.val)
}

func (g *RWValue[T]) Read(f func(*T)) {
	g.mux.RLock()
	defer g.mux.RUnlock()
	f(&g.val)
}
