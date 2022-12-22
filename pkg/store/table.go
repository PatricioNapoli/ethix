package store

import "sync"

type Table[T any] struct {
	Data map[string]T
	mut  *sync.Mutex
}

func NewTable[T any]() *Table[T] {
	return &Table[T]{
		Data: map[string]T{},
		mut:  &sync.Mutex{},
	}
}

func (t *Table[T]) Init() {
	t.mut = &sync.Mutex{}
}

func (t *Table[T]) Lock() {
	t.mut.Lock()
}

func (t *Table[T]) Unlock() {
	t.mut.Unlock()
}

func (t *Table[T]) Set(key string, value T) {
	t.mut.Lock()
	defer t.mut.Unlock()

	t.Data[key] = value
}

func (t *Table[T]) Get(key string) T {
	t.mut.Lock()
	defer t.mut.Unlock()

	return t.Data[key]
}
