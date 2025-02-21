package util

import "iter"

type RB[T any] struct {
	data    []T
	pointer int
	length  int
}

func NewRB[T any](cap int) *RB[T] {
	return &RB[T]{
		data:    make([]T, cap, cap),
		pointer: 0,
	}
}

func (rb *RB[T]) Append(newData T) {
	rb.data[rb.pointer] = newData
	rb.pointer = (rb.pointer + 1) % len(rb.data)
	if rb.length < len(rb.data) {
		rb.length += 1
	}
}

func (rb *RB[T]) Slice(s, e int, def T) iter.Seq[T] {
	return func(yield func(T) bool) {
		if s < 0 || e < 0 || s >= e {
			return
		}
		ptr := rb.pointer
		cap := len(rb.data)
		if s >= cap {
			return
		}
		end := min(e, rb.length)
		empty := e - end
		cur := ptr - end
		if cur < 0 {
			cur = cap + cur
		}
		for range e - s {
			if empty > 0 {
				if !yield(def) {
					return
				}
				empty -= 1
				continue
			}
			if !yield(rb.data[cur]) {
				return
			}
			cur = (cur + 1) % cap
		}
	}
}
